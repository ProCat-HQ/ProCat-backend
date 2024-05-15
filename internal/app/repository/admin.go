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
	query := fmt.Sprintf(`SELECT d.id, d.time_start, d.time_end, d.method, o.address, o.latitude, o.longitude, d.order_id, d.deliveryman_id
								 FROM %s d INNER JOIN %s o ON d.order_id = o.id WHERE o.status = $1`, deliveriesTable, ordersTable)

	queryDeliveryMan := fmt.Sprintf(`SELECT id, car_capacity, working_hours_start, working_hours_end  FROM %s`, deliverymenTable)

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
	query := fmt.Sprintf(`UPDATE %s SET deliveryman_id = $1 WHERE id = $2`, deliveriesTable)
	for point, deliverymanId := range answerMap {
		if _, err := a.db.Exec(query, deliverymanId, point.DeliveryId); err != nil {
			return err
		}
	}
	return nil
}

func (a *AdminPostgres) GetActualDeliveries() ([]model.DeliveryAddress, []model.Id, error) {
	query := fmt.Sprintf(`SELECT d.id, o.address, o.latitude, o.longitude, d.deliveryman_id
								 FROM %s d INNER JOIN %s o ON d.order_id = o.id WHERE o.status = $1`, deliveriesTable, ordersTable)
	queryDeliverymen := fmt.Sprintf(`SELECT id FROM %s`, deliverymenTable)
	var deliveries []model.DeliveryAddress
	if err := a.db.Select(&deliveries, query, "to_delivery"); err != nil {
		return nil, nil, err
	}
	var ids []model.Id
	err := a.db.Select(&ids, queryDeliverymen)
	if err != nil {
		return nil, nil, err
	}
	return deliveries, ids, nil
}

func (a *AdminPostgres) ChangeDeliveryman(delivery int, deliverymanId int) error {
	query := fmt.Sprintf(`UPDATE %s SET deliveryman_id = $1 WHERE id = $2`, deliveriesTable)
	_, err := a.db.Exec(query, deliverymanId, delivery)
	if err != nil {
		return err
	}
	return nil
}
