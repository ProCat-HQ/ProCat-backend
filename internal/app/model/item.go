package model

import (
	"database/sql"
	"time"
)

type PieceOfItem struct {
	Id          int            `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	Description sql.NullString `json:"description" db:"description"`
	Price       int            `json:"price" db:"price"`
	IsInStock   bool           `json:"isInStock" db:"is_in_stock"`
	Images      []string       `json:"images"`
	CategoryId  sql.NullInt32  `json:"categoryId" db:"category_id"`
}

type PieceOfItemToRes struct {
	Id          int      `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Description *string  `json:"description" db:"description"`
	Price       int      `json:"price" db:"price"`
	IsInStock   bool     `json:"isInStock" db:"is_in_stock"`
	Images      []string `json:"images"`
	CategoryId  *int     `json:"categoryId" db:"category_id"`
}

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
