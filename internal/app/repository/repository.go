package repository

import "github.com/jmoiron/sqlx"

type User interface {
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
	return &Repository{}
}
