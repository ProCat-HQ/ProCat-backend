package model

import "time"

type User struct {
	Id                   int       `json:"id"`
	FullName             string    `json:"fullName"`
	Email                string    `json:"email"`
	PhoneNumber          string    `json:"phoneNumber"`
	IdentificationNumber string    `json:"identificationNumber"`
	Password             string    `json:"password"`
	IsConfirmed          bool      `json:"isConfirmed"`
	Role                 string    `json:"role"`
	CreatedAt            time.Time `json:"createdAt"`
}

type DeliveryMan struct {
	Id                int       `json:"id"`
	CarCapacity       string    `json:"carCapacity"`
	WorkingHoursStart time.Time `json:"workingHoursStart"`
	WorkingHoursEnd   time.Time `json:"workingHoursEnd"`
	CarId             string    `json:"carId"`
	UserId            int
}

type SignUpInput struct {
	FullName    string `json:"fullName" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
}
