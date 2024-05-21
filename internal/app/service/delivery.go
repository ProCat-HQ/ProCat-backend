package service

import (
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"github.com/procat-hq/procat-backend/internal/routing"
	"strconv"
)

type DeliveryService struct {
	repo repository.Delivery
}

func NewDeliveryService(repo repository.Delivery) *DeliveryService {
	return &DeliveryService{repo: repo}
}

func (s *DeliveryService) GetDeliveriesForDeliveryman(userId int) (*model.MapRequest, map[model.LatLon]model.Point, error) {
	deliverymanId, err := s.repo.GetDeliverymanId(userId)
	if err != nil {
		return nil, nil, err
	}

	deliveries, err := s.repo.GetDeliveriesOrdersForDeliveryman(deliverymanId)
	if err != nil {
		return nil, nil, err
	}
	var sources []int
	var targets []int
	var points []model.LatLon
	mapDeliveriesId := make(map[model.LatLon]model.Point)
	for i, delivery := range deliveries {
		sources = append(sources, i)
		targets = append(targets, i)
		lat, err := strconv.ParseFloat(delivery.Latitude, 64)
		if err != nil {
			return nil, nil, err
		}
		lon, err := strconv.ParseFloat(delivery.Longitude, 64)
		if err != nil {
			return nil, nil, err
		}
		point := model.LatLon{
			Latitude:  lat,
			Longitude: lon,
		}
		points = append(points, point)
		mapDeliveriesId[point] = model.Point{
			DeliveryId: delivery.Id,
			Latitude:   lat,
			Longitude:  lon,
			Address:    delivery.Address,
		}
	}
	req := &model.MapRequest{
		Points:  points,
		Sources: sources,
		Targets: targets,
	}

	return req, mapDeliveriesId, nil
}

func (s *DeliveryService) GetAllDeliveries(statuses []string, limit string, page string, idStr string) ([]model.DeliveryWithOrder, int, error) {
	lim, err := strconv.Atoi(limit)
	if err != nil {
		return nil, 0, err
	}
	pag, err := strconv.Atoi(page)
	if err != nil {
		return nil, 0, err
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, 0, err
	}
	deliveries, count, err := s.repo.GetAllDeliveries(statuses, lim, lim*pag, id)
	if err != nil {
		return nil, 0, err
	}
	return deliveries, count, nil
}

func (s *DeliveryService) GetDelivery(idStr string) (model.DeliveryWithOrder, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return model.DeliveryWithOrder{}, err
	}
	delivery, err := s.repo.GetDelivery(id)
	if err != nil {
		return model.DeliveryWithOrder{}, err
	}
	return delivery, nil
}

func (s *DeliveryService) ChangeDeliveryStatus(idStr string, newStatus string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	err = s.repo.ChangeDeliveryStatus(id, newStatus)
	if err != nil {
		return err
	}
	return nil
}

func (s *DeliveryService) CreateRoute(requestBody model.MapRequest, responseFromApi model.Api2GisResponse,
	mapDeliveriesPoint map[model.LatLon]model.Point, userId int) ([]model.Point, error) {
	deliverymanId, err := s.repo.GetDeliverymanId(userId)
	if err != nil {
		return nil, err
	}
	result, err := routing.GetRoute(responseFromApi)
	if err != nil {
		return nil, err
	}
	var response []model.Point
	for _, i := range result.OptimalRoute {
		response = append(response, mapDeliveriesPoint[requestBody.Points[i]])
	}
	err = s.repo.InsertRoute(response, deliverymanId)
	return response, nil
}

func (s *DeliveryService) CheckRoute(userId int) ([]model.Point, error) {
	deliverymanId, err := s.repo.GetDeliverymanId(userId)
	if err != nil {
		return nil, err
	}
	route, err := s.repo.GetRoute(deliverymanId)
	if err != nil {
		return nil, err
	}
	return route, nil
}
