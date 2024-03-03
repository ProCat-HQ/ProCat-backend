package repository

type Authorization interface {
}

type Items interface {
}

type Repository struct {
	Authorization
	Items
}

func NewRepository() *Repository {
	return &Repository{}
}
