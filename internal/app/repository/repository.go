package repository

type Authorization interface {
}

type Items interface {
}

type Service struct {
	Authorization
	Items
}

func NewService() *Service {
	return &Service{}
}
