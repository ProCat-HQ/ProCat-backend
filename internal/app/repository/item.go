package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
)

type ItemPostgres struct {
	db *sqlx.DB
}

func NewItemPostgres(db *sqlx.DB) *ItemPostgres {
	return &ItemPostgres{db: db}
}

func (r *ItemPostgres) GetAllItems(limit, offset, categoryId int, stock bool) ([]model.PieceOfItem, error) {
	query := fmt.Sprintf("SELECT id, name, description, price, is_in_stock, category_id FROM %s OFFSET $1 LIMIT $2", itemsTable)

	var items []model.PieceOfItem

	err := r.db.Select(&items, query, offset, limit)
	if err != nil {
		return nil, err
	}

	return items, nil
}
