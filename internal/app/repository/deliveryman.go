package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
)

type DeliverymanPostgres struct {
	db *sqlx.DB
}

func NewDeliverymanPostgres(db *sqlx.DB) *DeliverymanPostgres {
	return &DeliverymanPostgres{db: db}
}

func (r *DeliverymanPostgres) GetAllDeliverymen(limit int, offset int) ([]model.DeliveryManInfoDB, error) {
	query := fmt.Sprintf(`SELECT d.id, d.car_capacity, coalesce(cast(d.working_hours_start as varchar), '') as working_hours_start,
								   coalesce(cast(d.working_hours_end as varchar), '') as working_hours_end,
								   d.car_id, u.fullname, u.email, u.phone_number
									FROM %s d join %s u on d.user_id = u.id OFFSET $1 LIMIT $2`, deliverymanTable, usersTable)
	var deliverymen []model.DeliveryManInfoDB
	err := r.db.Select(&deliverymen, query, offset, limit)
	if err != nil {
		return nil, err
	}
	return deliverymen, nil
}

func (r *DeliverymanPostgres) GetDeliveryman(deliveryId int) (*model.DeliveryManInfoCreate, error) {
	query := fmt.Sprintf(`SELECT d.car_capacity, coalesce(cast(d.working_hours_start as varchar), '') as working_hours_start,
								   coalesce(cast(d.working_hours_end as varchar), '') as working_hours_end, d.car_id
								FROM %s d
								JOIN %s ds on ds.deliveryman_id = d.id
								WHERE ds.id = $1`, deliverymanTable, deliveriesTable)
	var deliveryman model.DeliveryManInfoCreate
	err := r.db.Get(&deliveryman, query, deliveryId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &deliveryman, nil
}

func (r *DeliverymanPostgres) CreateDeliveryman(newDeliveryman model.DeliveryManInfoCreate, userId int) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (car_capacity, working_hours_start, working_hours_end, car_id, user_id) 
								VALUES ($1, $2, $3, $4, $5) 
								RETURNING id`, deliverymanTable)
	var id int
	err := r.db.Get(&id, query, newDeliveryman.CarCapacity, newDeliveryman.WorkingHoursStart, newDeliveryman.WorkingHoursEnd, newDeliveryman.CarId, userId)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *DeliverymanPostgres) ChangeDeliverymanData(newData model.DeliveryManInfoCreate, deliverymanId int) error {
	if newData.CarId != "" {
		query := fmt.Sprintf(`UPDATE %s
									SET car_id = $1
									WHERE id = $2`, deliverymanTable)
		_, err := r.db.Exec(query, newData.CarId, deliverymanId)
		if err != nil {
			return err
		}
	}
	if newData.WorkingHoursEnd != "" {
		query := fmt.Sprintf(`UPDATE %s
									SET working_hours_end = $1
									WHERE id = $2`, deliverymanTable)
		_, err := r.db.Exec(query, newData.WorkingHoursEnd, deliverymanId)
		if err != nil {
			return err
		}
	}
	if newData.WorkingHoursStart != "" {
		query := fmt.Sprintf(`UPDATE %s
									SET working_hours_start = $1
									WHERE id = $2`, deliverymanTable)
		_, err := r.db.Exec(query, newData.WorkingHoursStart, deliverymanId)
		if err != nil {
			return err
		}
	}
	if newData.CarCapacity != "" {
		query := fmt.Sprintf(`UPDATE %s
									SET car_capacity = $1
									WHERE id = $2`, deliverymanTable)
		_, err := r.db.Exec(query, newData.CarCapacity, deliverymanId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *DeliverymanPostgres) DeleteDeliveryman(deliverymanId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, deliverymanTable)
	_, err := r.db.Exec(query, deliverymanId)
	if err != nil {
		return err
	}
	return nil
}
