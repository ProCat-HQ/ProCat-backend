package model

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Payload any    `json:"payload"`
}

type Point struct {
	Address    string  `json:"address" db:"address"`
	Latitude   float64 `json:"latitude" db:"latitude"`
	Longitude  float64 `json:"longitude" db:"longitude"`
	DeliveryId int     `json:"deliveryId" db:"id"`
}

type Id struct {
	Id int `json:"id" db:"id"`
}
