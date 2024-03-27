package model

import "time"

type Chat struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	IsSolved     bool   `json:"isSolved"`
	FirstUserId  int
	SecondUserId int
	OrderId      int
}

type Message struct {
	Id        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	UserId    int
	ChatId    int
}

type MessageImage struct {
	Id        int    `json:"id"`
	Image     string `json:"image"`
	MessageId int
}
