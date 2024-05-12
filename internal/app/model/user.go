package model

import (
	"time"
)

type User struct {
	Id                   int       `json:"id" db:"id"`
	FullName             string    `json:"fullName" db:"fullname"`
	Email                string    `json:"email" db:"email"`
	PhoneNumber          string    `json:"phoneNumber" db:"phone_number"`
	IdentificationNumber string    `json:"identificationNumber" db:"identification_number"`
	IsConfirmed          bool      `json:"isConfirmed" db:"is_confirmed"`
	Role                 string    `json:"role" db:"role"`
	CreatedAt            time.Time `json:"createdAt" db:"created_at"`
}

type DeliveryMan struct {
	Id                int       `json:"id" db:"id"`
	CarCapacity       string    `json:"carCapacity" db:"car_capacity"`
	WorkingHoursStart time.Time `json:"workingHoursStart" db:"working_hours_start"`
	WorkingHoursEnd   time.Time `json:"workingHoursEnd" db:"working_hours_end"`
	CarId             string    `json:"carId" db:"car_id"`
	UserId            int       `db:"user_id"`
}

type RefreshSession struct {
	Id           int    `json:"id" db:"id"`
	RefreshToken string `json:"refreshToken" db:"refresh_token"`
	Fingerprint  string `json:"fingerprint" db:"fingerprint"`
	UserId       int    `json:"userId" db:"user_id"`
}

// DeliveriesForDeliveryMan переместить
type DeliveriesForDeliveryMan struct {
	DeliverymanId int     `json:"deliverymanId"`
	Deliveries    []Point `json:"deliveries"`
}

type SignUpInput struct {
	FullName    string `json:"fullName" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type AccessTokenClaimsExtension struct {
	UserId   int    `json:"userId"`
	UserRole string `json:"userRole"`
}

type RefreshTokenClaimsExtension struct {
	UserId int `json:"userId"`
}
