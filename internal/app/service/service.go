package service

import "github.com/procat-hq/procat-backend/internal/app/repository"

type Authorization interface {
}

type Items interface {
}

type Service struct {
	Authorization
	Items
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
