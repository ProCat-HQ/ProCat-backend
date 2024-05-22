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
	OrderId    int          `json:"orderId"`
	TotalPrice int          `json:"totalPrice"`
	Deposit    int          `json:"deposit"`
	Items      []ItemCheque `json:"items"`
}

type Order struct {
	Id                int              `json:"id" db:"id"`
	Status            string           `json:"status" db:"status"`
	TotalPrice        int              `json:"totalPrice" db:"total_price"`
	Deposit           int              `json:"deposit" db:"deposit"`
	RentalPeriodStart time.Time        `json:"rentalPeriodStart" db:"rental_period_start"`
	RentalPeriodEnd   time.Time        `json:"rentalPeriodEnd" db:"rental_period_end"`
	Address           string           `json:"address" db:"address"`
	Latitude          string           `json:"latitude" db:"latitude"`
	Longitude         string           `json:"longitude" db:"longitude"`
	CompanyName       string           `json:"companyName" db:"company_name"`
	CreatedAt         time.Time        `json:"createdAt" db:"created_at"`
	UserId            int              `json:"userId" db:"user_id"`
	Items             []OrderSmallItem `json:"items"`
}

type OrderSmallItem struct {
	Id           int    `json:"id" db:"item_id"`
	Name         string `json:"name" db:"name"`
	Price        int    `json:"price" db:"price"`
	PriceDeposit int    `json:"priceDeposit" db:"price_deposit"`
	Count        int    `json:"count" db:"count"`
	Image        string `json:"image" db:"image"`
}

type DeliveryWithOrder struct {
	Id            int       `json:"id" db:"id"`
	TimeStart     time.Time `json:"timeStart" db:"time_start"`
	TimeEnd       time.Time `json:"timeEnd" db:"time_end"`
	Method        string    `json:"method" db:"method"`
	DeliveryManId int       `json:"deliveryManId" db:"deliveryman_id"`
	SmallOrder    `json:"order"`
}

type SmallOrder struct {
	OrderId    int    `json:"id" db:"order_id"`
	Status     string `json:"status" db:"status"`
	TotalPrice int    `json:"totalPrice" db:"total_price"`
	Address    string `json:"address" db:"address"`
	Latitude   string `json:"latitude" db:"latitude"`
	Longitude  string `json:"longitude" db:"longitude"`
}

type OrderSmall struct {
	Id         int    `json:"id" db:"id"`
	Status     string `json:"status" db:"status"`
	TotalPrice int    `json:"totalPrice" db:"total_price"`
	Address    string `json:"address" db:"address"`
	Latitude   string `json:"latitude" db:"latitude"`
	Longitude  string `json:"longitude" db:"longitude"`
}

type DeliveryAddress struct {
	Point
	DeliverymanId int `json:"deliverymanId" db:"deliveryman_id"`
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
	Id        int       `json:"id" db:"id"`
	Paid      int       `json:"paid" db:"paid"`
	Method    string    `json:"method" db:"method"`
	Price     int       `json:"price" db:"price"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type WaitingHoursForRouting struct {
	Id    int
	Start int
	End   int
}
