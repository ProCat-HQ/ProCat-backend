package model

type Payment struct {
	Id      int    `json:"id"`
	IsPaid  bool   `json:"isPaid"`
	Method  string `json:"method"`
	Price   int    `json:"price"`
	OrderId int
}
