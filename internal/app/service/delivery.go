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

func (s *DeliveryService) GetAllDeliveries(statuses []string, limit string, page string, idStr string) ([]model.DeliveryFullInfo, int, error) {
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
	payload := make([]model.DeliveryFullInfo, len(deliveries))
	for i, delivery := range deliveries {
		makeDeliveryNote(&payload[i], delivery)
	}
	return payload, count, nil
}

func (s *DeliveryService) GetDelivery(idStr string) (*model.DeliveryFullInfo, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	delivery, err := s.repo.GetDelivery(id)
	if err != nil {
		return nil, err
	}
	var payload model.DeliveryFullInfo
	makeDeliveryNote(&payload, *delivery)
	return &payload, nil
}

func makeDeliveryNote(note *model.DeliveryFullInfo, info model.OrderAndDeliveryInfo) {
	note.Id = info.Id
	note.TimeStart = info.TimeStart
	note.TimeEnd = info.TimeEnd
	note.Method = info.Method
	note.DeliveryManId = info.DeliveryManId
	note.Order.Id = info.OrderId
	note.Order.Status = info.Status
	note.Order.TotalPrice = info.TotalPrice
	note.Order.Address = info.Address
	note.Order.Latitude = info.Latitude
	note.Order.Longitude = info.Longitude
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
