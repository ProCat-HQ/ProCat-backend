package model

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Payload any    `json:"payload"`
}

type Point struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	DeliveryId int     `json:"deliveryId"`
}
