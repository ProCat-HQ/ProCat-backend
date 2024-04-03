package service

import (
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"github.com/procat-hq/procat-backend/internal/routing"
	"github.com/sirupsen/logrus"
)

type AdminService struct {
	repo repository.Admin
}

func NewAdminService(repo repository.Admin) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) MakeClustering() error {
	// логика вся
	deliveries, deliveryMen, err := s.repo.GetDeliveries()
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	answer, err := routing.ClusterOrders(deliveries, deliveryMen)
	if err != nil {
		return err
	}
	err = s.repo.SetDeliveries(answer)
	if err != nil {
		return err
	}
	return nil
}
