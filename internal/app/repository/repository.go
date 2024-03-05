package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
}

type Items interface {
}

type Repository struct {
	Authorization
	Items
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
