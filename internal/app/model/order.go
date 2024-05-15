package model

import (
	"database/sql"
	"time"
)

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

type OrderAndDeliveryInfo struct {
	Id            int       `json:"id" db:"id"`
	TimeStart     time.Time `json:"timeStart" db:"time_start"`
	TimeEnd       time.Time `json:"timeEnd" db:"time_end"`
	Method        string    `json:"method" db:"method"`
	DeliveryManId int       `json:"deliveryManId" db:"deliveryman_id"`
	OrderId       int       `json:"orderId" db:"order_id"`
	Status        string    `json:"status" db:"status"`
	TotalPrice    int       `json:"totalPrice" db:"total_price"`
	Address       string    `json:"address" db:"address"`
	Latitude      string    `json:"latitude" db:"latitude"`
	Longitude     string    `json:"longitude" db:"longitude"`
}

type OrderSmall struct {
	Id         int    `json:"id" db:"id"`
	Status     string `json:"status" db:"status"`
	TotalPrice int    `json:"totalPrice" db:"total_price"`
	Address    string `json:"address" db:"address"`
	Latitude   string `json:"latitude" db:"latitude"`
	Longitude  string `json:"longitude" db:"longitude"`
}

type DeliveryFullInfo struct {
	Id            int        `json:"id" db:"id"`
	TimeStart     time.Time  `json:"timeStart" db:"time_start"`
	TimeEnd       time.Time  `json:"timeEnd" db:"time_end"`
	Method        string     `json:"method" db:"method"`
	DeliveryManId int        `json:"deliveryManId" db:"deliveryman_id"`
	Order         OrderSmall `json:"order"`
}

type Delivery struct {
	Id            int       `json:"id"`
	TimeStart     time.Time `json:"timeStart"`
	TimeEnd       time.Time `json:"timeEnd"`
	Method        string    `json:"method"`
	OrderId       int
	DeliveryManId int
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
