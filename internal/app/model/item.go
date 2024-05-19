package model

import (
	"time"
)

type CartItem struct {
	Id    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Price int    `json:"price" db:"price"`
	Count int    `json:"count" db:"count"`
	Image string `json:"image" db:"image"`
}

type PieceOfItem struct {
	Id           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Description  string `json:"description" db:"description"`
	Price        int    `json:"price" db:"price"`
	IsInStock    bool   `json:"isInStock" db:"is_in_stock"`
	CategoryId   int    `json:"categoryId" db:"category_id"`
	CategoryName string `json:"categoryName" db:"category_name"`
	Image        string `json:"image" db:"image"`
}

type Item struct {
	Id           int          `json:"id" db:"id"`
	Name         string       `json:"name" db:"name"`
	Description  string       `json:"description" db:"description"`
	Price        int          `json:"price" db:"price"`
	IsInStock    bool         `json:"isInStock" db:"is_in_stock"`
	CategoryId   int          `json:"categoryId" db:"category_id"`
	CategoryName string       `json:"categoryName" db:"category_name"`
	Info         []Info       `json:"info"`
	Images       []ItemImage  `json:"images"`
	ItemStores   []ItemStores `json:"itemStores"`
}

type Info struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

type ItemImage struct {
	Id    int    `json:"id" db:"id"`
	Image string `json:"image" db:"image"`
}

type ItemStores struct {
	Id                int       `json:"id" db:"id"`
	InStockNumber     int       `json:"inStockNumber" db:"in_stock_number"`
	Name              string    `json:"name" db:"name"`
	Address           string    `json:"address" db:"address"`
	WorkingHoursStart time.Time `json:"workingHoursStart" db:"working_hours_start"`
	WorkingHoursEnd   time.Time `json:"workingHoursEnd" db:"working_hours_end"`
}

type Store struct {
	Id                int       `json:"id" db:"id"`
	Name              string    `json:"name" db:"name" binding:"required"`
	Address           string    `json:"address" db:"address" binding:"required"`
	Latitude          float64   `json:"latitude" db:"latitude"`
	Longitude         float64   `json:"longitude" db:"longitude"`
	WorkingHoursStart time.Time `json:"workingHoursStart" db:"working_hours_start" binding:"required"`
	WorkingHoursEnd   time.Time `json:"workingHoursEnd" db:"working_hours_end" binding:"required"`
}

type StoreFromDB struct {
	Id                int       `json:"id" db:"id"`
	Name              string    `json:"name" db:"name" binding:"required"`
	Address           string    `json:"address" db:"address" binding:"required"`
	Latitude          string    `json:"latitude" db:"latitude"`
	Longitude         string    `json:"longitude" db:"longitude"`
	WorkingHoursStart time.Time `json:"workingHoursStart" db:"working_hours_start" binding:"required"`
	WorkingHoursEnd   time.Time `json:"workingHoursEnd" db:"working_hours_end" binding:"required"`
}

type StoreCreation struct {
	Name              string `json:"name" binding:"required"`
	Address           string `json:"address" binding:"required"`
	WorkingHoursStart string `json:"workingHoursStart" binding:"required"`
	WorkingHoursEnd   string `json:"workingHoursEnd" binding:"required"`
}

type StoreChange struct {
	Name              string `json:"name"`
	Address           string `json:"address"`
	WorkingHoursStart string `json:"workingHoursStart"`
	WorkingHoursEnd   string `json:"workingHoursEnd"`
}

type StoreChangeDB struct {
	Name              string
	Address           string
	Latitude          float64
	Longitude         float64
	WorkingHoursStart *time.Time
	WorkingHoursEnd   *time.Time
}

type ChangeStock struct {
	StoreId       int `json:"storeId" binding:"required"`
	InStockNumber int `json:"inStockNumber" binding:"min=0"`
}
