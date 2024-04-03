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
}

type Delivery interface {
}

type Admin interface {
	MakeClustering() (string, error)
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
		Admin: NewAdminService(repos.Admin),
		User:  NewUserService(repos.User),
	}
}
