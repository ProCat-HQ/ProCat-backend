package service

import (
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"mime/multipart"
)

type User interface {
	CreateUser(user model.SignUpInput) (int, error)
	GetUserByCredentials(phoneNumber, password string) (model.User, error)
	GenerateTokens(user model.User, fingerprint string) (string, string, error)
	ParseAccessToken(accessToken string) (*model.AccessTokenClaimsExtension, error)
	ParseRefreshToken(refreshToken string) (*model.RefreshTokenClaimsExtension, error)
	LogoutUser(refreshToken string, userId int) (int, error)
	RegenerateTokens(userId int, refreshToken, fingerprint string) (string, string, error)

	GetAllUsers(limit, page, role, isConfirmed string) (int, []model.User, error)
	GetUserById(userId int) (model.User, error)
	DeleteUserById(userId int) error
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
	AddItemsToCart(userId, itemId, count int) error
	DeleteItemFromCart(userId, itemId int) error
	GetCartItems(userId int) ([]model.CartItem, error)
}

type Order interface {
	CreateOrder(userId int, order model.OrderCreationWithTime) (model.OrderCheque, error)
}

type Subscription interface {
}

type Notification interface {
}

type Category interface {
}

type Item interface {
	GetAllItems(limit, page, search, categoryId, stock string) (int, []model.PieceOfItem, error)
	GetItem(itemId string) (model.Item, error)
	CreateItem(name, description, price, priceDeposit, categoryId string, files []*multipart.FileHeader) (int, error)
	DeleteItem(itemId int) error
	ChangeItem(itemId int, name, description, price, priceDeposit, categoryId *string) error
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
		Cart:     NewCartService(repos.Cart),
		Order:    NewOrderService(repos.Order),
	}
}
