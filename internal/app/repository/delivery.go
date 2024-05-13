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
	err := r.db.Select(&deliveries, query, "to_delivery", deliverymanId)
	if err != nil {
		return nil, err
	}
	return deliveries, nil
}
