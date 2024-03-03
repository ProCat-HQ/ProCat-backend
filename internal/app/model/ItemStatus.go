package model

type ItemStatus struct {
	Id            int  `json:"id"`
	IsInStock     bool `json:"isInStock"`
	InStockNumber int  `json:"inStockNumber"`
	ItemId        int
}
