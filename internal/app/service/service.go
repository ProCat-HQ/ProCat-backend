package service

import (
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
)

type User interface {
	CreateUser(user model.SignUpInput) (int, error)
	GetUserByCredentials(phoneNumber, password string) (model.User, error)
	GenerateTokens(user model.User, fingerprint string) (string, string, error)
	ParseAccessToken(accessToken string) (*model.AccessTokenClaimsExtension, error)
	ParseRefreshToken(refreshToken string) (*model.RefreshTokenClaimsExtension, error)
	LogoutUser(refreshToken string, userId int) (int, error)
	RegenerateTokens(userId int, refreshToken, fingerprint string) (string, string, error)
}

type Verification interface {
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
		User:     NewUserService(repos.User),
		Item:     NewItemService(repos.Item),
		Admin:    NewAdminService(repos.Admin),
		Delivery: NewDeliveryService(repos.Delivery),
	}
}
