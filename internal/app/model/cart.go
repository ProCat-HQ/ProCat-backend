package model

type Cart struct {
	Id     int
	UserId int
}

type CartItem struct {
	Id          int
	ItemsNumber int
	CartId      int
	ItemId      int
}
