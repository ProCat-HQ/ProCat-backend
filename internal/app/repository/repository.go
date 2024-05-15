package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"time"
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

type Deliveryman interface {
	GetAllDeliverymen(limit int, offset int) ([]model.DeliveryManInfoDB, int, error)
	GetDeliveryman(userId int) (model.DeliveryManInfoCreate, error)
	CreateDeliveryman(newDeliveryman model.DeliveryManInfoCreate, userId int) (int, error)
	ChangeDeliverymanData(newData model.DeliveryManInfoCreate, deliverymanId int) error
	DeleteDeliveryman(deliverymanId int) error
}

type Delivery interface {
	GetDeliverymanId(userId int) (int, error)
	GetDeliveriesOrdersForDeliveryman(deliverymanId int) ([]model.DeliveryAndOrder, error)
	GetAllDeliveries(statuses []string, limit int, offset int, id int) ([]model.DeliveryWithOrder, int, error)
	GetDelivery(id int) (model.DeliveryWithOrder, error)
	ChangeDeliveryStatus(id int, newStatus string) error
}

type Admin interface {
	GetDeliveries() ([]model.DeliveryAndOrder, []model.DeliveryMan, error)
	SetDeliveries(map[model.Point]int) error
	GetDeliveriesToSort() (int, []model.DeliveriesForDeliveryMan, error)
	ChangeDeliveryman(delivery int, deliverymanId int) error
}

type Cart interface {
	GetUsersCartId(userId int) (int, error)
	AddItemToCart(cartId, itemId, count int) error
	DeleteItemFromCart(cartId, itemId int) error
	GetCartItems(cartId int) ([]model.CartItem, error)
}

type Order interface {
	GetUserById(userId int) (model.User, error)
	GetUsersCartId(userId int) (int, error)
	GetTotalCartPrices(cartId int) (int, int, error)
	GetItemCheque(cartId int) ([]model.ItemCheque, error)
	CreateOrder(status string, deposit bool, rpStart, rpEnd time.Time,
		address string, lat, lon float64, companyName string, userId int,
		deliveryMethod string, tStart, tEnd time.Time) (model.OrderCheque, error)
}

type Subscription interface {
}

type Notification interface {
}

type Category interface {
}

type Item interface {
	GetCategoryChildren(categoryId int) ([]int, error)
	GetAllItems(limit, offset, categoryId int, stock bool, search string) (int, []model.PieceOfItem, error)
	GetItem(itemId int) (model.Item, error)
	CreateItem(name, description string, price, priceDeposit, categoryId int) (int, error)
	SaveFilenames(itemId int, filenames []string) error
	DeleteItem(itemId int) error
	ChangeItem(itemId int, name, description, price, priceDeposit, categoryId *string) error
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
		Cart:        NewCartPostgres(db),
		Order:       NewOrderPostgres(db),
	}
}
