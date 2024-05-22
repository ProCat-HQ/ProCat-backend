package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"time"
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

func (r *DeliveryPostgres) GetWorkingHours(deliverymanId int) (int, error) {
	var start time.Time
	query := fmt.Sprintf(`SELECT working_hours_start
								FROM %s 
								WHERE id = $1`, deliverymenTable)
	err := r.db.Get(&start, query, deliverymanId)
	if err != nil {
		return 0, err
	}
	return start.Hour(), nil
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

func (r *DeliveryPostgres) GetRoute(deliverymanId int) ([]model.Point, error) {
	queryRoute := fmt.Sprintf(`SELECT COALESCE(
               (SELECT id FROM %s WHERE deliveryman_id = $1 LIMIT 1), -1) AS result`, routesTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var routeId int
	err = tx.Get(&routeId, queryRoute, deliverymanId)
	if err != nil {
		return nil, err
	}
	if routeId == -1 {
		return nil, nil
	}
	queryStatuses := fmt.Sprintf(`SELECT status 
										FROM %s 
											JOIN %s 
												ON deliveries.order_id = orders.id 
													   AND deliveryman_id = $1
										WHERE status IN ('%s', '%s')`,
		ordersTable, deliveriesTable, model.ReadyToDelivery, model.Delivering)
	var statuses []struct {
		Status string `db:"status"`
	}
	err = r.db.Select(&statuses, queryStatuses, deliverymanId)
	if err != nil {
		return nil, err
	}
	queryCountCoordinates := fmt.Sprintf(`SELECT count(*) FROM %s where route_id = $1`, coordinatesTable)
	var count int
	err = tx.Get(&count, queryCountCoordinates, routeId)
	if err != nil {
		return nil, err
	}
	queryChangeStatus := fmt.Sprintf(`UPDATE %s AS o SET status=$1
											FROM %s AS d JOIN %s c on d.id = c.delivery_id
											WHERE o.id = d.order_id AND d.deliveryman_id = $2`,
		ordersTable, deliveriesTable, coordinatesTable)
	queryDeleteRoute := fmt.Sprintf(`DELETE FROM routes where id = $1`)
	if len(statuses) != count {
		_, err = tx.Exec(queryDeleteRoute, routeId)
		if err != nil {
			return nil, err
		}
		_, err = tx.Exec(queryChangeStatus, model.ReadyToDelivery, deliverymanId)
		if err != nil {
			return nil, err
		}
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	for _, status := range statuses {
		if status.Status == model.ReadyToDelivery {
			_, err = tx.Exec(queryDeleteRoute, routeId)
			if err != nil {
				return nil, err
			}
			_, err = tx.Exec(queryChangeStatus, model.ReadyToDelivery, deliverymanId)
			if err != nil {
				return nil, err
			}
			err = tx.Commit()
			if err != nil {
				return nil, err
			}
			return nil, nil
		}
	}
	query := fmt.Sprintf(`SELECT o.latitude, o.longitude, address, d.id
								FROM %s o
									JOIN %s d
										ON o.id = d.order_id
									JOIN %s c
										ON d.id = c.delivery_id
								WHERE route_id = $1
								ORDER BY sequence_number`, ordersTable, deliveriesTable, coordinatesTable)
	var route []model.Point
	err = tx.Select(&route, query, routeId)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return route, nil
}

func (r *DeliveryPostgres) InsertRoute(route []model.Point, deliverymanId int) error {
	queryRouteId := fmt.Sprintf(`INSERT INTO %s (deliveryman_id) VALUES ($1) RETURNING id`, routesTable)
	var routeId int

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.Get(&routeId, queryRouteId, deliverymanId)
	if err != nil {
		return err
	}
	queryChangeStatus := fmt.Sprintf(`UPDATE %s AS o SET status=$1
											FROM %s AS d WHERE o.id = d.order_id
											AND d.id=$2`, ordersTable, deliveriesTable)

	queryPoints := fmt.Sprintf(`INSERT INTO %s (sequence_number, delivery_id, route_id) 
								VALUES ($1, $2, $3)`, coordinatesTable)

	stmtStatus, err := tx.Preparex(queryChangeStatus)
	if err != nil {
		return err
	}
	defer stmtStatus.Close()
	stmtPoints, err := tx.Preparex(queryPoints)
	if err != nil {
		return err
	}
	defer stmtPoints.Close()

	for i, point := range route {
		_, err = stmtStatus.Exec(model.Delivering, point.DeliveryId)
		if err != nil {
			return err
		}

		_, err = stmtPoints.Exec(i, point.DeliveryId, routeId)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *DeliveryPostgres) GetStore(storeId int) (model.Point, error) {
	query := fmt.Sprintf(`SELECT address, latitude, longitude FROM %s WHERE id = $1`, storesTable)
	var store model.Point
	err := r.db.Get(&store, query, storeId)
	return store, err
}
