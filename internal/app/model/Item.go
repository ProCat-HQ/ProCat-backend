package model

type Item struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	CategoryId  int
	//TODO: SimilarToItems
}

type ItemStatus struct {
	Id            int    `json:"id"`
	IsInStock     bool   `json:"isInStock"`
	InStockNumber int    `json:"inStockNumber"`
	Address       string `json:"address"`
	ItemId        int
}

type ItemImage struct {
	Id     int    `json:"id"`
	Image  string `json:"image"`
	ItemId int
}
