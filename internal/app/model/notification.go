package model

import "time"

type Notification struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsViewed    bool      `json:"isViewed"`
	CreatedAt   time.Time `json:"createdAt"`
	UserId      int
}
