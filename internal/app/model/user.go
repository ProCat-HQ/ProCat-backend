package model

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
