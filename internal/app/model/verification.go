package model

import "time"

type Verification struct {
	Id       int       `json:"id"`
	Code     string    `json:"code"`
	Type     string    `json:"type"`
	Value    string    `json:"value"`
	LifeTime time.Time `json:"lifeTime"`
	UserId   int
}
