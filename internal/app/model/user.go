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

type UserPassword struct {
	Id                   int       `json:"id" db:"id"`
	FullName             string    `json:"fullName" db:"fullname"`
	Email                string    `json:"email" db:"email"`
	PhoneNumber          string    `json:"phoneNumber" db:"phone_number"`
	IdentificationNumber string    `json:"identificationNumber" db:"identification_number"`
	PasswordHash         string    `json:"passwordHash" db:"password_hash"`
	IsConfirmed          bool      `json:"isConfirmed" db:"is_confirmed"`
	Role                 string    `json:"role" db:"role"`
	CreatedAt            time.Time `json:"createdAt" db:"created_at"`
}

type RefreshSession struct {
	Id           int    `json:"id" db:"id"`
	RefreshToken string `json:"refreshToken" db:"refresh_token"`
	Fingerprint  string `json:"fingerprint" db:"fingerprint"`
	UserId       int    `json:"userId" db:"user_id"`
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
