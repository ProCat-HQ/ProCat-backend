package model

import "time"

type Item struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	IsInStock   bool   `json:"isInStock"`
	CategoryId  int
	//TODO: SimilarToItems
}

type ItemStores struct {
	Id            int `json:"id"`
	InStockNumber int `json:"inStockNumber"`
	StoreId       int
	ItemId        int
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

type ItemImage struct {
	Id     int    `json:"id"`
	Image  string `json:"image"`
	ItemId int
}
