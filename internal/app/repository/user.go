package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
)

type UserPostgres struct {
	db *sqlx.DB
}

func (r *UserPostgres) GetUser(phoneNumber, password string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id, role FROM %s WHERE phone_number=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, phoneNumber, password)
	return user, err
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) CreateUser(user model.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (fullname, phone_number, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.FullName, user.PhoneNumber, user.Password)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
