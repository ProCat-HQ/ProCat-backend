package model

import "time"

type User struct {
	Id                   int       `json:"id" db:"id"`
	FullName             string    `json:"fullName"`
	Email                string    `json:"email"`
	PhoneNumber          string    `json:"phoneNumber"`
	IdentificationNumber string    `json:"identificationNumber"`
	Password             string    `json:"password"`
	IsConfirmed          bool      `json:"isConfirmed"`
	Role                 string    `json:"role" db:"role"`
	CreatedAt            time.Time `json:"createdAt"`
}

type DeliveryMan struct {
	Id                int       `json:"id" db:"id"`
	CarCapacity       string    `json:"carCapacity" db:"car_capacity"`
	WorkingHoursStart time.Time `json:"workingHoursStart" db:"working_hours_start"`
	WorkingHoursEnd   time.Time `json:"workingHoursEnd" db:"working_hours_end"`
	CarId             string    `json:"carId" db:"car_id"`
	UserId            int       `db:"user_id"`
}

type SignUpInput struct {
	FullName    string `json:"fullName" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type TokenClaimsExtension struct {
	UserId   int    `json:"userId"`
	UserRole string `json:"userRole"`
}
