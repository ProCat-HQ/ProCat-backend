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

	GetUserWithPasswordById(userId int) (model.UserPassword, error)
	ChangeFullName(userId int, fullName string) error
	ChangeIdentificationNumber(userId int, identificationNumber string) error
	ChangePassword(userId int, passwordHash string) error
	ChangePhoneNumber(userId int, phoneNumber, passwordHash string) error
	ChangeEmail(userId int, email string) error
	ChangeUserRole(userId int, role string) error
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
	ChangeDeliveryman(deliveryId int, deliverymanId int) error
}

type Cart interface {
	GetUsersCartId(userId int) (int, error)
	AddItemToCart(cartId, itemId, count int) error
	DeleteItemFromCart(cartId, itemId, count int) error
	GetCartItems(cartId int) ([]model.CartItem, error)
}

type Order interface {
	GetAllOrders(limit, offset, userId int, statuses []string) (int, []model.Order, error)
	GetOrder(orderId int) (model.Order, error)
	GetUserById(userId int) (model.User, error)
	GetUsersCartId(userId int) (int, error)
	GetTotalCartPrices(cartId int) (int, int, error)
	GetItemCheque(cartId int) ([]model.ItemCheque, error)
	CreateOrder(status string, deposit bool, rpStart, rpEnd time.Time,
		address string, lat, lon float64, companyName string, userId int,
		deliveryMethod string, tStart, tEnd time.Time, rentPeriodDays int) (model.OrderCheque, error)
	ChangeOrderStatus(orderId int, status string) error
	GetPaymentsForOrder(orderId int) ([]model.Payment, error)
	ChangePaymentStatus(paymentId, paid int, method string) error
}

type Subscription interface {
	GetUserSubscriptions(userId int, limit, offset int) (int, []model.Subscription, error)
	CreateSubscription(userId, itemId int) error
	DeleteSubscription(userId, subId int) error
}

type Notification interface {
	GetUsersNotification(userId int) ([]model.Notification, error)
	CreateNotification(userId int, title, description string) (int, error)
	ReadAndGetNotification(notificationId int) (model.Notification, error)
	GetNotificationUserId(notificationId int) (int, error)
	DeleteNotification(notificationId int) error
}

type Category interface {
	CreateCategory(categoryParentId int, name string) (int, error)
	ChangeCategory(categoryId int, name string) error
	GetCategoriesForParent(categoryParentId int) ([]model.Category, error)
	DeleteCategory(categoryId int) error
	GetCategoryRoute(categoryId int) ([]model.Category, error)
}

type Item interface {
	GetCategoryChildren(categoryId int) ([]int, error)
	GetAllItems(limit, offset, categoryId int, stock bool, search string) (int, []model.PieceOfItem, error)
	GetItem(itemId int) (model.Item, error)
	CreateItem(name, description string, price, priceDeposit, categoryId int) (int, error)
	SaveFilenames(itemId int, filenames []string) error
	DeleteItem(itemId int) error
	ChangeItem(itemId int, name, description, price, priceDeposit, categoryId *string) error

	ChangeStockOfItem(itemId, storeId, inStockNumber int) error

	AddInfos(itemId int, info model.ItemInfoCreation) error
	ChangeInfos(itemId int, info model.ItemInfoChange) error
	DeleteInfos(itemId int, ids []int) error

	DeleteImages(itemId int, ids []int) ([]string, error)
}

type Store interface {
	CreateStore(store model.Store) (int, error)
	GetAllStores() ([]model.StoreFromDB, error)
	ChangeStore(storeId int, store model.StoreChangeDB) error
	DeleteStore(storeId int) error
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
		User:         NewUserPostgres(db),
		Item:         NewItemPostgres(db),
		Admin:        NewAdminPostgres(db),
		Deliveryman:  NewDeliverymanPostgres(db),
		Delivery:     NewDeliveryPostgres(db),
		Cart:         NewCartPostgres(db),
		Order:        NewOrderPostgres(db),
		Store:        NewStorePostgres(db),
		Category:     NewCategoryPostgres(db),
		Notification: NewNotificationPostgres(db),
		Subscription: NewSubscriptionPostgres(db),
	}
}
