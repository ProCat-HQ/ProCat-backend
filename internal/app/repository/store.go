package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
)

type StorePostgres struct {
	db *sqlx.DB
}

func NewStorePostgres(db *sqlx.DB) *StorePostgres {
	return &StorePostgres{db: db}
}

func (r *StorePostgres) CreateStore(store model.Store) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (name, address, latitude, longitude,
                				 working_hours_start, working_hours_end)
								 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, storesTable)

	var id int
	err := r.db.Get(&id, query, store.Name, store.Address, store.Latitude, store.Longitude,
		store.WorkingHoursStart, store.WorkingHoursEnd)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *StorePostgres) GetAllStores() ([]model.StoreFromDB, error) {
	query := fmt.Sprintf(`SELECT id, id, name, address, COALESCE(latitude, '') AS latitude,
       							COALESCE(longitude, '') AS longitude, working_hours_start, working_hours_end
								FROM %s`, storesTable)

	var stores []model.StoreFromDB
	err := r.db.Select(&stores, query)
	if err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *StorePostgres) ChangeStore(storeId int, store model.StoreChangeDB) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if store.Name != "" {
		query := fmt.Sprintf(`UPDATE %s SET name=$1 WHERE id=$2`, storesTable)
		_, err = tx.Exec(query, store.Name, storeId)
		if err != nil {
			return err
		}
	}
	if store.Address != "" {
		query := fmt.Sprintf(`UPDATE %s SET address=$1, latitude=$2, longitude=$3 WHERE id=$4`, storesTable)
		_, err = tx.Exec(query, store.Address, store.Latitude, store.Longitude, storeId)
		if err != nil {
			return err
		}
	}
	if store.WorkingHoursStart != nil {
		query := fmt.Sprintf(`UPDATE %s SET working_hours_start=$1 WHERE id=$2`, storesTable)
		_, err = tx.Exec(query, *store.WorkingHoursStart, storeId)
		if err != nil {
			return err
		}
	}
	if store.WorkingHoursEnd != nil {
		query := fmt.Sprintf(`UPDATE %s SET working_hours_end=$1 WHERE id=$2`, storesTable)
		_, err = tx.Exec(query, *store.WorkingHoursEnd, storeId)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *StorePostgres) DeleteStore(storeId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, storesTable)
	_, err := r.db.Exec(query, storeId)
	return err
}
