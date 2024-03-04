package model

import "time"

type Order struct {
	Id                int       `json:"id"`
	Status            string    `json:"status"`
	TotalPrice        int       `json:"totalPrice"`
	RentalPeriodStart time.Time `json:"rentalPeriodStart"`
	RentalPeriodEnd   time.Time `json:"rentalPeriodEnd"`
	Address           string    `json:"address"`
	CompanyName       string    `json:"companyName"`
	Contact           string    `json:"contact"`
	UserId            int
}

type Delivery struct {
	Id            int       `json:"id"`
	Time          time.Time `json:"time"`
	Method        string    `json:"method"`
	OrderId       int
	DeliveryManId int
}

type Payment struct {
	Id      int    `json:"id"`
	IsPaid  bool   `json:"isPaid"`
	Method  string `json:"method"`
	Price   int    `json:"price"`
	OrderId int
}

type OrderItem struct {
	Id          int
	ItemsNumber int
	OrderId     int
	ItemId      int
}
