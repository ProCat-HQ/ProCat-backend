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

func (r *AdminPostgres) GetDeliveries() ([]model.DeliveryAndOrder, []model.DeliveryMan, error) {
	query := fmt.Sprintf(`SELECT d.id, d.time_start, d.time_end, d.method, o.address, o.latitude, o.longitude, d.order_id, d.deliveryman_id
								 FROM %s d INNER JOIN %s o ON d.order_id = o.id WHERE o.status = $1`, deliveriesTable, ordersTable)

	queryDeliveryMan := fmt.Sprintf(`SELECT id, car_capacity, working_hours_start, working_hours_end  FROM %s`, deliverymenTable)

	var deliveries []model.DeliveryAndOrder
	var deliverymen []model.DeliveryMan
	if err := r.db.Select(&deliveries, query, model.Accepted); err != nil {
		return nil, nil, err
	}
	if err := r.db.Select(&deliverymen, queryDeliveryMan); err != nil {
		return nil, nil, err
	}
	return deliveries, deliverymen, nil

}

func (r *AdminPostgres) SetDeliveries(answerMap map[model.Point]int) error {
	query := fmt.Sprintf(`UPDATE %s SET deliveryman_id = $1 WHERE id = $2`, deliveriesTable)
	for point, deliverymanId := range answerMap {
		if _, err := r.db.Exec(query, deliverymanId, point.DeliveryId); err != nil {
			return err
		}
	}
	return nil
}

func (r *AdminPostgres) GetDeliveriesToSort() (int, []model.DeliveriesForDeliveryMan, error) {
	queryDeliverymen := fmt.Sprintf(`SELECT d.deliveryman_id
											FROM %s d
											JOIN %s o ON d.order_id = o.id
											WHERE d.deliveryman_id IS NOT NULL
											AND o.status = $1
											GROUP BY d.deliveryman_id
											ORDER BY deliveryman_id`, deliveriesTable, ordersTable)

	queryDeliveries := fmt.Sprintf(`SELECT o.address, o.latitude, o.longitude, d.id
										   FROM %s d
										   JOIN %s o ON d.order_id = o.id
										   WHERE d.deliveryman_id=$1
										   AND o.status = $2`, deliveriesTable, ordersTable)

	var deliverymenId []int
	err := r.db.Select(&deliverymenId, queryDeliverymen, model.Accepted)
	if err != nil {
		return 0, nil, err
	}
	response := make([]model.DeliveriesForDeliveryMan, 0)

	for _, deliverymanId := range deliverymenId {
		var deliveries []model.Point
		err = r.db.Select(&deliveries, queryDeliveries, deliverymanId, model.Accepted)
		if err != nil {
			return 0, nil, err
		}
		res := model.DeliveriesForDeliveryMan{
			DeliverymanId: deliverymanId,
			Deliveries:    deliveries,
		}
		response = append(response, res)
	}
	return len(response), response, nil
}

func (r *AdminPostgres) ChangeDeliveryman(delivery int, deliverymanId int) error {
	query := fmt.Sprintf(`UPDATE %s SET deliveryman_id = $1 WHERE id = $2`, deliveriesTable)
	_, err := r.db.Exec(query, deliverymanId, delivery)
	if err != nil {
		return err
	}
	return nil
}
