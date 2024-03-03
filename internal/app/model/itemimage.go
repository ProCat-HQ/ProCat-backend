package model

type ItemImage struct {
	Id     int    `json:"id"`
	Image  string `json:"image"`
	ItemId int
}
