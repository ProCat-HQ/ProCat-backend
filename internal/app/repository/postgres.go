package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable              = "users"
	deliverymanTable        = "deliverymen"
	verificationsTable      = "verifications"
	ordersTable             = "orders"
	deliveriesTable         = "deliveries"
	paymentsTable           = "payments"
	categoriesTable         = "categories"
	itemsTable              = "items"
	storesTable             = "stores"
	itemsStoresTable        = "item_stores"
	itemsImagesTable        = "item_images"
	ordersItemsTable        = "orders_items"
	infosTable              = "infos"
	cartsTable              = "carts"
	cartsItemsTable         = "carts_items"
	subscriptionsTable      = "subscriptions"
	subscriptionsItemsTable = "subscriptions_items"
	notificationsTable      = "notifications"
	chatsTable              = "chats"
	messagesTable           = "messages"
	messageImagesTable      = "message_images"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
