package model

import (
	"database/sql"
	"time"
)

type OrderCreation struct {
	RentalPeriodStart string `json:"rentalPeriodStart" binding:"required"`
	RentalPeriodEnd   string `json:"rentalPeriodEnd" binding:"required"`
	Address           string `json:"address" binding:"required"`
	CompanyName       string `json:"companyName"`
	DeliveryMethod    string `json:"deliveryMethod" binding:"required"`
	TimeStart         string `json:"deliveryTimeStart" binding:"required"`
	TimeEnd           string `json:"deliveryTimeEnd" binding:"required"`
}

type OrderCreationWithTime struct {
	RentalPeriodStart time.Time `json:"rentalPeriodStart"`
	RentalPeriodEnd   time.Time `json:"rentalPeriodEnd"`
	Address           string    `json:"address"`
	CompanyName       string    `json:"companyName"`
	DeliveryMethod    string    `json:"deliveryMethod"`
	TimeStart         time.Time `json:"deliveryTimeStart"`
	TimeEnd           time.Time `json:"deliveryTimeEnd"`
}

type ItemCheque struct {
	Name         string `json:"name" db:"name"`
	Count        int    `json:"count" db:"count"`
	Price        int    `json:"price" db:"price"`
	PriceDeposit int    `json:"priceDeposit" db:"price_deposit"`
}

type OrderCheque struct {
	OrderId      int          `json:"orderId"`
	TotalPrice   int          `json:"totalPrice"`
	TotalDeposit int          `json:"totalDeposit"`
	Items        []ItemCheque `json:"items"`
}

type Order struct {
	Id                int       `json:"id"`
	Status            string    `json:"status"`
	TotalPrice        int       `json:"totalPrice"`
	Deposit           int       `json:"deposit"`
	RentalPeriodStart time.Time `json:"rentalPeriodStart"`
	RentalPeriodEnd   time.Time `json:"rentalPeriodEnd"`
	Address           string    `json:"address"`
	Latitude          string    `json:"latitude"`
	Longitude         string    `json:"longitude"`
	CompanyName       string    `json:"companyName"`
	CreatedAt         time.Time `json:"createdAt"`
	UserId            int
}

type Delivery struct {
	Id            int       `json:"id"`
	TimeStart     time.Time `json:"timeStart"`
	TimeEnd       time.Time `json:"timeEnd"`
	Method        string    `json:"method"`
	OrderId       int
	DeliveryManId int
}

type LatLon struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type MapRequest struct {
	Points  []LatLon `json:"points"`
	Sources []int    `json:"sources"`
	Targets []int    `json:"targets"`
}

type DeliveryAndOrder struct {
	Id            int           `json:"id" db:"id"`
	TimeStart     time.Time     `json:"timeStart" db:"time_start"`
	TimeEnd       time.Time     `json:"timeEnd" db:"time_end"`
	Method        string        `json:"method" db:"method"`
	Address       string        `json:"address" db:"address"`
	Latitude      string        `json:"latitude" db:"latitude"`
	Longitude     string        `json:"longitude" db:"longitude"`
	OrderId       int           `db:"order_id"`
	DeliveryManId sql.NullInt64 `db:"deliveryman_id"`
}

type Payment struct {
	Id        int       `json:"id"`
	Paid      int       `json:"paid"`
	Method    string    `json:"method"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	OrderId   int
}

type OrderItem struct {
	Id          int
	ItemsNumber int
	OrderId     int
	ItemId      int
}
