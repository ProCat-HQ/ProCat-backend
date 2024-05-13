package service

import (
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
)

type User interface {
	CreateUser(user model.SignUpInput) (int, error)
	GenerateToken(phoneNumber, password string) (string, error)
	ParseToken(accessToken string) (*model.TokenClaimsExtension, error)
}

type Verification interface {
}

type Deliveryman interface {
	GetAllDeliverymen(limit string, page string) ([]model.DeliveryManInfoDB, error)
	GetDeliveryman(deliveryId string) (*model.DeliveryManInfoCreate, error)
	CreateDeliveryman(newDeliveryman model.DeliveryManInfoCreate, userId string) (int, error)
	ChangeDeliverymanData(newData model.DeliveryManInfoCreate, deliverymanId string) error
	DeleteDeliveryman(deliverymanId string) error
}

type Delivery interface {
	GetDeliveriesForDeliveryman(userId int) (*model.MapRequest, error)
}

type Admin interface {
	MakeClustering() ([]model.DeliveriesForDeliveryMan, error)
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
	GetAllItems(limit, page, categoryId, stock string) ([]model.PieceOfItem, error)
}

type Store interface {
}

type Service struct {
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

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User:        NewUserService(repos.User),
		Item:        NewItemService(repos.Item),
		Admin:       NewAdminService(repos.Admin),
		Deliveryman: NewDeliverymanService(repos.Deliveryman),
		Delivery:    NewDeliveryService(repos.Delivery),
	}
}
