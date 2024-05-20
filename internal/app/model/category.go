package model

type Category struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	ParentId int    `json:"parentId" db:"parent_id"`
}
