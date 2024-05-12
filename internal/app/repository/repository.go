package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
)

type User interface {
	CreateUser(user model.SignUpInput) (int, error)
	GetUser(phoneNumber, password string) (model.User, error)
	GetUserById(userId int) (model.User, error)
	GetRefreshSessions(userId int) ([]model.RefreshSession, error)
	GetRefreshSession(refreshToken string, userId int) (model.RefreshSession, error)
	WipeRefreshSessionsWithFingerprint(fingerprint string, userId int) error
	WipeRefreshSessions(userId int) error
	SaveSessionData(refreshToken, fingerprint string, userId int) error
	DeleteUserRefreshSession(refreshToken string, userId int) (int, error)

	GetAllUsers(limit, offset int, role, isConfirmed string) (int, []model.User, error)
	DeleteUserById(userId int) error
}

type Verification interface {
}

type Delivery interface {
	GetDeliverymanId(userId int) (int, error)
	GetDeliveriesOrdersForDeliveryman(deliverymanId int) ([]model.DeliveryAndOrder, error)
}

type Admin interface {
	GetDeliveries() ([]model.DeliveryAndOrder, []model.DeliveryMan, error)
	SetDeliveries(map[model.Point]int) error
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
		User:     NewUserPostgres(db),
		Item:     NewItemPostgres(db),
		Admin:    NewAdminPostgres(db),
		Delivery: NewDeliveryPostgres(db),
	}
}
