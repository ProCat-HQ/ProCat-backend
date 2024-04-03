package model

import "time"

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

type DeliveryAndOrder struct {
	Id            int       `json:"id"`
	TimeStart     time.Time `json:"timeStart"`
	TimeEnd       time.Time `json:"timeEnd"`
	Method        string    `json:"method"`
	Address       string    `json:"address"`
	Latitude      string    `json:"latitude"`
	Longitude     string    `json:"longitude"`
	OrderId       int
	DeliveryManId int
}

type Payment struct {
	Id        int       `json:"id"`
	IsPaid    bool      `json:"isPaid"`
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
