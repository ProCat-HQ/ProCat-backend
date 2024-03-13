package service

import (
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
)

type Authorization interface {
}

type Routes interface {
	SortPoints(points model.RouteList) (model.RouteList, error)
}

type Service struct {
	Authorization
	Routes
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
