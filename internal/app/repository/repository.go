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
	GetAllDeliverymen(limit int, offset int) ([]model.DeliveryManInfoDB, int, error)
	GetDeliveryman(deliverymanId int) (*model.DeliveryManInfoCreate, error)
	CreateDeliveryman(newDeliveryman model.DeliveryManInfoCreate, userId int) (int, error)
	ChangeDeliverymanData(newData model.DeliveryManInfoCreate, deliverymanId int) error
	DeleteDeliveryman(deliverymanId int) error
}

type Delivery interface {
	GetDeliverymanId(userId int) (int, error)
	GetDeliveriesOrdersForDeliveryman(deliverymanId int) ([]model.DeliveryAndOrder, error)
	GetAllDeliveries(statuses []string, limit int, offset int, id int) ([]model.OrderAndDeliveryInfo, int, error)
	GetDelivery(id int) (*model.OrderAndDeliveryInfo, error)
	ChangeDeliveryStatus(id int, newStatus string) error
}

type Admin interface {
	GetDeliveries() ([]model.DeliveryAndOrder, []model.DeliveryMan, error)
	SetDeliveries(map[model.Point]int) error
	GetActualDeliveries() ([]model.DeliveryAddress, []model.Id, error)
	ChangeDeliveryman(delivery int, deliverymanId int) error
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
	GetAllItems(limit, offset, categoryId int, stock bool) ([]model.PieceOfItem, error)
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
		User:        NewUserPostgres(db),
		Item:        NewItemPostgres(db),
		Admin:       NewAdminPostgres(db),
		Deliveryman: NewDeliverymanPostgres(db),
		Delivery:    NewDeliveryPostgres(db),
	}
}
