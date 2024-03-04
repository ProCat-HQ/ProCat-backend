package model

import "time"

type User struct {
	Id                   int    `json:"id"`
	FullName             string `json:"fullName"`
	Email                string `json:"email"`
	PhoneNumber          string `json:"phoneNumber"`
	IdentificationNumber string `json:"identificationNumber"`
	Password             string `json:"password"`
	IsConfirmed          bool   `json:"isConfirmed"`
	Role                 string `json:"role"`
}

type DeliveryMan struct {
	Id                int           `json:"id"`
	CarCapacity       string        `json:"carCapacity"`
	WorkingHoursStart time.Duration `json:"workingHoursStart"`
	WorkingHoursEnd   time.Duration `json:"workingHoursEnd"`
	CarId             string        `json:"carId"`
	UserId            int
}
