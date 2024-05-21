package model

import "time"

type Notification struct {
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	IsViewed    bool      `json:"isViewed" db:"is_viewed"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type NotificationCreate struct {
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description" binding:"required"`
}
