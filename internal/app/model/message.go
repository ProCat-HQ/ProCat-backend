package model

type Message struct {
	Id     int    `json:"id"`
	Text   string `json:"text"`
	UserId int
	ChatId int
}
