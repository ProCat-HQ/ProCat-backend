package model

import "github.com/gin-gonic/gin"

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Payload gin.H  `json:"payload"`
}

type Point struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	DeliveryId int     `json:"deliveryId"`
}
