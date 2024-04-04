package service

import (
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"strconv"
)

type DeliveryService struct {
	repo repository.Delivery
}

func NewDeliveryService(repo repository.Delivery) *DeliveryService {
	return &DeliveryService{repo: repo}
}

func (s *DeliveryService) GetDeliveriesForDeliveryman(userId int) (*model.MapRequest, error) {
	deliverymanId, err := s.repo.GetDeliverymanId(userId)
	if err != nil {
		return nil, err
	}

	deliveries, err := s.repo.GetDeliveriesOrdersForDeliveryman(deliverymanId)
	if err != nil {
		return nil, err
	}
	var sources []int
	var targets []int
	var points []model.LatLon

	for i, delivery := range deliveries {
		sources = append(sources, i)
		targets = append(targets, i)
		lat, err := strconv.ParseFloat(delivery.Latitude, 64)
		if err != nil {
			return nil, err
		}
		lon, err := strconv.ParseFloat(delivery.Longitude, 64)
		if err != nil {
			return nil, err
		}
		points = append(points, model.LatLon{
			Latitude:  lat,
			Longitude: lon,
		})
	}
	req := &model.MapRequest{
		Points:  points,
		Sources: sources,
		Targets: targets,
	}

	return req, nil
}
