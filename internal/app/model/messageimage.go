package model

type MessageImage struct {
	Id        int    `json:"id"`
	Image     string `json:"image"`
	MessageId int
}
