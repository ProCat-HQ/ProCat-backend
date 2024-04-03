package service

import (
	"encoding/json"
	"fmt"
	"github.com/procat-hq/procat-backend/internal/app/model"
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

func (s *AdminService) MakeClustering() (string, error) {
	deliveries, deliveryMen, err := s.repo.GetDeliveries()
	if err != nil {
		logrus.Error(err.Error())
		return "", err
	}
	answer, err := routing.ClusterOrders(deliveries, deliveryMen)
	if err != nil {
		return "", err
	}
	err = s.repo.SetDeliveries(answer)
	if err != nil {
		return "", err
	}
	var result []model.DeliveriesForDeliveryMan
	for _, man := range deliveryMen {
		result = append(result, model.DeliveriesForDeliveryMan{DeliveryManId: man.Id, Deliveries: make([]model.Point, 0)})
	}
	for i, i2 := range answer {
		for j, j2 := range result {
			if j2.DeliveryManId == i2 {
				result[j].Deliveries = append(result[j].Deliveries, i)
			}
		}
	}
	j, _ := json.Marshal(result)
	fmt.Print(string(j))
	return string(j), nil
}
