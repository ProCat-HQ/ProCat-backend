package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
)

type User interface {
	CreateUser(user model.User) (int, error)
	GetUser(phoneNumber, password string) (model.User, error)
}

type Verification interface {
}

type Deliveryman interface {
}

type Delivery interface {
}

type Admin interface {
}

type Cart interface {
}

type Order interface {
}

type Subscription interface {
}

type Notification interface {
}

type Category interface {
}

type Item interface {
}

type Store interface {
}

type Repository struct {
	User
	Verification
	Deliveryman
	Delivery
	Admin
	Cart
	Order
	Subscription
	Notification
	Category
	Item
	Store
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserPostgres(db),
	}
}