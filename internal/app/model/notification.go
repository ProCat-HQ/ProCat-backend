package model

type Notification struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsViewed    bool   `json:"isViewed"`
	UserId      int
}
