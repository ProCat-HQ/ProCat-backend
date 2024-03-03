package model

import "time"

type DeliveryMan struct {
	Id                int           `json:"id"`
	CarCapacity       string        `json:"carCapacity"`
	WorkingHoursStart time.Duration `json:"workingHoursStart"`
	WorkingHoursEnd   time.Duration `json:"workingHoursEnd"`
	CarId             string        `json:"carId"`
	UserId            int
}
