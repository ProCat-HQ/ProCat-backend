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

	CheckPassword(password string, userId int) (bool, error)
	ChangeFullName(userId int, fullName string) error
	ChangeIdentificationNumber(userId int, identificationNumber string) error
	ChangePassword(userId int, password string) error
	ChangePhoneNumber(userId int, phoneNumber, password string) error
	ChangeEmail(userId int, email string) error
	ChangeUserRole(userId int, role string) error
}

type Verification interface {
}

type Deliveryman interface {
	GetAllDeliverymen(limit string, page string) ([]model.DeliveryManInfoDB, int, error)
	GetDeliveryman(userId string) (model.DeliveryManInfoCreate, error)
	CreateDeliveryman(newDeliveryman model.DeliveryManInfoCreate, userId string) (int, error)
	ChangeDeliverymanData(newData model.DeliveryManInfoCreate, deliverymanId string) error
	DeleteDeliveryman(deliverymanId string) error
}

type Delivery interface {
	GetDeliveriesForDeliveryman(userId int) (*model.MapRequest, error)
	GetAllDeliveries(statuses []string, limit string, page string, idStr string) ([]model.DeliveryWithOrder, int, error)
	GetDelivery(idStr string) (model.DeliveryWithOrder, error)
	ChangeDeliveryStatus(id string, newStatus string) error
}

type Admin interface {
	MakeClustering() ([]model.DeliveriesForDeliveryMan, error)
	GetDeliveriesToSort() (int, []model.DeliveriesForDeliveryMan, error)
	ChangeDeliveryman(deliveryId int, deliverymanId int) error
}

type Cart interface {
	AddItemsToCart(userId, itemId, count int) error
	DeleteItemFromCart(userId, itemId, count int) error
	GetCartItems(userId int) ([]model.CartItem, error)
}

type Order interface {
	GetAllOrders(limit, page, userId int, statuses []string) (int, []model.Order, error)
	GetOrder(orderId int) (model.Order, error)
	CreateOrder(userId int, order model.OrderCreationWithTime) (model.OrderCheque, error)
	ChangeOrderStatus(orderId int, status string) error
	GetPaymentsForOrder(orderId int) ([]model.Payment, error)
	ChangePaymentStatus(paymentId, paid int, method string) error
}

type Subscription interface {
	GetUserSubscriptions(userId int, limit, page string) (int, []model.Subscription, error)
	CreateSubscription(userId, itemId int) error
	DeleteSubscription(userId, subId int) error
}

type Notification interface {
	GetUsersNotification(userId int) ([]model.Notification, error)
	CreateNotification(userId int, title, description string) (int, error)
	ReadAndGetNotification(userId int, notificationId int) (model.Notification, error)
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
	GetAllItems(limit, page, search, categoryId, stock string) (int, []model.PieceOfItem, error)
	GetItem(itemId string) (model.Item, error)
	CreateItem(name, description, price, priceDeposit, categoryId string, files []*multipart.FileHeader) (int, error)
	DeleteItem(itemId int) error
	ChangeItem(itemId int, name, description, price, priceDeposit, categoryId *string) error

	ChangeStockOfItem(itemId, storeId, inStockNumber int) error

	AddInfos(itemId int, info model.ItemInfoCreation) error
	ChangeInfos(itemId int, info model.ItemInfoChange) error
	DeleteInfos(itemId int, ids []int) error
	AddImages(itemId int, files []*multipart.FileHeader) error
	DeleteImages(itemId int, ids []int) error
}

type Store interface {
	CreateStore(store model.Store) (int, error)
	GetAllStores() ([]model.StoreFromDB, error)
	ChangeStore(storeId int, store model.StoreChange) error
	DeleteStore(storeId int) error
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
		User:         NewUserService(repos.User),
		Item:         NewItemService(repos.Item),
		Admin:        NewAdminService(repos.Admin),
		Delivery:     NewDeliveryService(repos.Delivery),
		Deliveryman:  NewDeliverymanService(repos.Deliveryman),
		Cart:         NewCartService(repos.Cart),
		Order:        NewOrderService(repos.Order),
		Store:        NewStoreService(repos.Store),
		Category:     NewCategoryService(repos.Category),
		Notification: NewNotificationService(repos.Notification),
		Subscription: NewSubscriptionService(repos.Subscription),
	}
}
