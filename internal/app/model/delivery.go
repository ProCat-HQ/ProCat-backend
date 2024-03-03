package model

import "time"

type Delivery struct {
	Id            int       `json:"id"`
	Time          time.Time `json:"time"`
	Method        string    `json:"method"`
	OrderId       int
	DeliveryManId int
}
