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
