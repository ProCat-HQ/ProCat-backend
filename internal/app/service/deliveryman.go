package service

import "github.com/procat-hq/procat-backend/internal/app/repository"

type DeliverymanService struct {
	repo repository.Deliveryman
}

func NewDeliverymanService(repo repository.Deliveryman) *DeliverymanService {
	return &DeliverymanService{repo: repo}
}
