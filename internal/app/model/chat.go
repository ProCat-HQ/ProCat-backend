package model

type Chat struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	IsSolved     bool   `json:"isSolved"`
	FirstUserId  int
	SecondUserId int
	OrderId      int
}
