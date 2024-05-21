package model

type Subscription struct {
	Id               int `json:"id" db:"id"`
	SubscriptionItem `json:"item"`
}

type SubscriptionItem struct {
	Id        int    `json:"id" db:"item_id"`
	Name      string `json:"name" db:"name"`
	Price     int    `json:"price" db:"price"`
	IsInStock bool   `json:"isInStock" db:"is_in_stock"`
	Image     string `json:"image" db:"image"`
}
