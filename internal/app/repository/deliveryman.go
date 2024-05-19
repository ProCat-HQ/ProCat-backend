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

func (r *DeliverymanPostgres) GetAllDeliverymen(limit int, offset int) ([]model.DeliveryManInfoDB, int, error) {
	queryForCount := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, deliverymenTable)
	var count int
	err := r.db.Get(&count, queryForCount)
	if err != nil {
		return nil, 0, err
	}
	query := fmt.Sprintf(`SELECT d.id, COALESCE(d.car_capacity, '') AS car_capacity, 
       							 COALESCE(CAST(d.working_hours_start AS VARCHAR), '') AS working_hours_start,
								 COALESCE(CAST(d.working_hours_end AS VARCHAR), '') AS working_hours_end,
								 COALESCE(d.car_id, '') AS car_id, u.fullname, COALESCE(u.email, '') AS email, u.phone_number
								 FROM %s d JOIN %s u ON d.user_id = u.id OFFSET $1 LIMIT $2`, deliverymenTable, usersTable)
	var deliverymen []model.DeliveryManInfoDB
	err = r.db.Select(&deliverymen, query, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	return deliverymen, count, nil
}

func (r *DeliverymanPostgres) GetDeliveryman(userId int) (model.DeliveryManInfoCreate, error) {
	query := fmt.Sprintf(`SELECT COALESCE(car_capacity, '') AS car_capacity,
							    COALESCE(CAST(working_hours_start AS VARCHAR), '') AS working_hours_start,
								COALESCE(CAST(working_hours_end AS VARCHAR), '') AS working_hours_end,
								COALESCE(car_id, '') AS car_id
								FROM %s d
								WHERE user_id = $1`, deliverymenTable)
	var deliveryman model.DeliveryManInfoCreate
	err := r.db.Get(&deliveryman, query, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.DeliveryManInfoCreate{}, nil
		}
		return model.DeliveryManInfoCreate{}, err
	}
	return deliveryman, nil
}

func (r *DeliverymanPostgres) CreateDeliveryman(newDeliveryman model.DeliveryManInfoCreate, userId int) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (car_capacity, working_hours_start, working_hours_end, car_id, user_id) 
								VALUES ($1, $2, $3, $4, $5) 
								RETURNING id`, deliverymenTable)

	queryChangeRole := fmt.Sprintf(`UPDATE %s SET role=$1 WHERE id=$2`, usersTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var id int
	err = tx.Get(&id, query, newDeliveryman.CarCapacity, newDeliveryman.WorkingHoursStart, newDeliveryman.WorkingHoursEnd, newDeliveryman.CarId, userId)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(queryChangeRole, model.DeliverymanRole, userId)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *DeliverymanPostgres) ChangeDeliverymanData(newData model.DeliveryManInfoCreate, deliverymanId int) error {
	if newData.CarId != "" {
		query := fmt.Sprintf(`UPDATE %s
									SET car_id = $1
									WHERE id = $2`, deliverymenTable)
		_, err := r.db.Exec(query, newData.CarId, deliverymanId)
		if err != nil {
			return err
		}
	}
	if newData.WorkingHoursEnd != "" {
		query := fmt.Sprintf(`UPDATE %s
									SET working_hours_end = $1
									WHERE id = $2`, deliverymenTable)
		_, err := r.db.Exec(query, newData.WorkingHoursEnd, deliverymanId)
		if err != nil {
			return err
		}
	}
	if newData.WorkingHoursStart != "" {
		query := fmt.Sprintf(`UPDATE %s
									SET working_hours_start = $1
									WHERE id = $2`, deliverymenTable)
		_, err := r.db.Exec(query, newData.WorkingHoursStart, deliverymanId)
		if err != nil {
			return err
		}
	}
	if newData.CarCapacity != "" {
		query := fmt.Sprintf(`UPDATE %s
									SET car_capacity = $1
									WHERE id = $2`, deliverymenTable)
		_, err := r.db.Exec(query, newData.CarCapacity, deliverymanId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *DeliverymanPostgres) DeleteDeliveryman(deliverymanId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1 RETURNING user_id`, deliverymenTable)
	queryChangeRole := fmt.Sprintf(`UPDATE %s SET role=$1 WHERE id=$2`, usersTable)

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userId int

	err = tx.Get(&userId, query, deliverymanId)
	if err != nil {
		return err
	}
	_, err = tx.Exec(queryChangeRole, model.UserRole, userId)
	if err != nil {
		return err
	}

	return tx.Commit()
}
