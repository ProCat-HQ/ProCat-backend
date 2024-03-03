package model

type Item struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	CategoryId  int
	//TODO: SimilarToItems
}
