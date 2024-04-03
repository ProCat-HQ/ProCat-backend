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

func (a *AdminPostgres) GetDeliveries() (*model.DeliveryAndOrder, error) {
	query := fmt.Sprintf("SELECT id, time_start, time_end, method, address, latitude," +
		"longitude, order_id, delivery_man_id FROM deliveries INNER JOIN orders ON deliveries.order_id == orders.id")
}
