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
	Id                int       `json:"id"`
	Name              string    `json:"name"`
	Address           string    `json:"address"`
	Latitude          string    `json:"latitude"`
	Longitude         string    `json:"longitude"`
	WorkingHoursStart time.Time `json:"workingHoursStart"`
	WorkingHoursEnd   time.Time `json:"workingHoursEnd"`
}
