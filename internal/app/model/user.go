package model

import (
	"time"
)

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

type SignUpInput struct {
	FullName    string `json:"fullName" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type TokenClaimsExtension struct {
	UserId   int    `json:"userId"`
	UserRole string `json:"userRole"`
}
