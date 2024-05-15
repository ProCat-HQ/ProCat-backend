package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
)

type DeliveryPostgres struct {
	db *sqlx.DB
}

func NewDeliveryPostgres(db *sqlx.DB) *DeliveryPostgres {
	return &DeliveryPostgres{db: db}
}

func (r *DeliveryPostgres) GetDeliverymanId(userId int) (int, error) {
	query := fmt.Sprintf(`SELECT id FROM %s WHERE user_id = $1`, deliverymenTable)
	row := r.db.QueryRow(query, userId)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *DeliveryPostgres) GetDeliveriesOrdersForDeliveryman(deliverymanId int) ([]model.DeliveryAndOrder, error) {
	query := fmt.Sprintf(`SELECT d.id, d.time_start, d.time_end, d.method, o.address, o.latitude, o.longitude,
       d.order_id, d.deliveryman_id FROM %s d INNER JOIN %s o ON d.order_id = o.id
                                                               WHERE o.status = $1 AND d.deliveryman_id = $2`, deliveriesTable, ordersTable)

	var deliveries []model.DeliveryAndOrder
	err := r.db.Select(&deliveries, query, model.ReadyToDelivery, deliverymanId)
	if err != nil {
		return nil, err
	}
	return deliveries, nil
}

func (r *DeliveryPostgres) GetAllDeliveries(statuses []string, limit int, offset int, id int) ([]model.DeliveryWithOrder, int, error) {
	query := fmt.Sprintf(`SELECT d.id, time_start, time_end, method, COALESCE(deliveryman_id, -1) AS deliveryman_id,
       							order_id, status, total_price, address, COALESCE(latitude, '') AS latitude,
       							COALESCE(longitude, '') AS longitude
								FROM %s d
								JOIN %s o ON o.id = d.order_id`, deliveriesTable, ordersTable)
	if len(statuses) > 0 || id >= 0 {
		query = query + ` WHERE`
	}
	if id >= 0 {
		query += ` deliveryman_id = $1`
		if len(statuses) > 0 {
			query += ` AND`
		}
	}
	if len(statuses) > 0 {
		query += ` status IN (`
		for i, status := range statuses {
			query = query + fmt.Sprintf(`'%s'`, status)
			if i != len(statuses)-1 {
				query = query + `, `
			} else {
				query += `)`
			}

		}
	}
	var count int
	queryCount := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, deliveriesTable)
	var err error
	if id == -1 {
		err = r.db.Get(&count, queryCount)
	} else {
		queryCount += ` WHERE deliveryman_id = $1`
		err = r.db.Get(&count, queryCount, id)
	}

	if err != nil {
		return nil, 0, err
	}

	var deliveries []model.DeliveryWithOrder
	if id == -1 {
		query = query + ` LIMIT $1 OFFSET $2`
		err = r.db.Select(&deliveries, query, limit, offset)
	} else {
		query = query + ` LIMIT $2 OFFSET $3`
		err = r.db.Select(&deliveries, query, id, limit, offset)
	}

	if err != nil {
		return nil, 0, err
	}

	return deliveries, count, nil
}

func (r *DeliveryPostgres) GetDelivery(id int) (model.DeliveryWithOrder, error) {
	query := fmt.Sprintf(`SELECT d.id, time_start, time_end, method, COALESCE(deliveryman_id, -1) AS deliveryman_id,
       							order_id, status, total_price, address, COALESCE(latitude, '') AS latitude,
       							COALESCE(longitude, '') AS longitude
								FROM %s d
								JOIN %s o ON o.id = d.order_id
								WHERE d.id = $1`, deliveriesTable, ordersTable)
	var delivery model.DeliveryWithOrder
	err := r.db.Get(&delivery, query, id)
	if err != nil {
		return model.DeliveryWithOrder{}, err
	}
	return delivery, nil
}

func (r *DeliveryPostgres) ChangeDeliveryStatus(id int, newStatus string) error {
	query := fmt.Sprintf(`UPDATE %s SET status = $1
								WHERE id
								IN (SELECT order_id FROM %s WHERE id = $2)`, ordersTable, deliveriesTable)
	_, err := r.db.Exec(query, newStatus, id)
	if err != nil {
		return err
	}
	return nil
}
