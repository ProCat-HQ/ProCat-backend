package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
)

type AdminPostgres struct {
	db *sqlx.DB
}

func NewAdminPostgres(db *sqlx.DB) *AdminPostgres {
	return &AdminPostgres{db: db}
}

func (a *AdminPostgres) GetDeliveries() ([]model.DeliveryAndOrder, []model.DeliveryMan, error) {
	query := fmt.Sprintf(`SELECT d.id, d.time_start, d.time_end, d.method, o.address, o.latitude, o.longitude, d.order_id, d.delivery_man_id
								 FROM %s d INNER JOIN %s o ON d.order_id = o.id WHERE o.status = $1`, deliveriesTable, ordersTable)
	queryDeliveryMan := fmt.Sprintf(`SELECT id, car_capacity, working_hours_start, working_hours_end  FROM delivery_men`)
	var deliveries []model.DeliveryAndOrder
	var deliverymen []model.DeliveryMan
	if err := a.db.Select(&deliveries, query, "to_delivery"); err != nil {
		return nil, nil, err
	}
	if err := a.db.Select(&deliverymen, queryDeliveryMan); err != nil {
		return nil, nil, err
	}
	return deliveries, deliverymen, nil

}

func (a *AdminPostgres) SetDeliveries(answerMap map[model.Point]int) error {
	query := fmt.Sprintf(`UPDATE %s SET delivery_man_id = $1 WHERE id = $2`, deliveriesTable)
	for point, deliverymanId := range answerMap {
		if _, err := a.db.Exec(query, deliverymanId, point.DeliveryId); err != nil {
			return err
		}
	}
	return nil
}
