package model

import "time"

type DeliveryMan struct {
	Id                int       `json:"id" db:"id"`
	CarCapacity       string    `json:"carCapacity" db:"car_capacity"`
	WorkingHoursStart time.Time `json:"workingHoursStart" db:"working_hours_start"`
	WorkingHoursEnd   time.Time `json:"workingHoursEnd" db:"working_hours_end"`
	CarId             string    `json:"carId" db:"car_id"`
	UserId            int       `db:"user_id"`
}

type DeliveryManInfoDB struct {
	Id                int    `json:"id" db:"id"`
	CarCapacity       string `json:"carCapacity" db:"car_capacity"`
	WorkingHoursStart string `json:"workingHoursStart" db:"working_hours_start"`
	WorkingHoursEnd   string `json:"workingHoursEnd" db:"working_hours_end"`
	CarId             string `json:"carId" db:"car_id"`
	FullName          string `json:"fullName" db:"fullname"`
	Email             string `json:"email" db:"email"`
	PhoneNumber       string `json:"phoneNumber" db:"phone_number"`
}

type DeliveryManInfoCreate struct {
	CarCapacity       string `json:"carCapacity" db:"car_capacity"`
	WorkingHoursStart string `json:"workingHoursStart" db:"working_hours_start"`
	WorkingHoursEnd   string `json:"workingHoursEnd" db:"working_hours_end"`
	CarId             string `json:"carId" db:"car_id"`
}

type DeliveriesForDeliveryMan struct {
	DeliverymanId int     `json:"deliverymanId"`
	Deliveries    []Point `json:"deliveries"`
}
