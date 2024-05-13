package service

import (
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"strconv"
)

type DeliverymanService struct {
	repo repository.Deliveryman
}

func NewDeliverymanService(repo repository.Deliveryman) *DeliverymanService {
	return &DeliverymanService{repo: repo}
}

func (s *DeliverymanService) GetAllDeliverymen(limit string, page string) ([]model.DeliveryManInfoDB, int, error) {
	lim, err := strconv.Atoi(limit)
	if err != nil {
		return nil, 0, err
	}
	pag, err := strconv.Atoi(page)
	if err != nil {
		return nil, 0, err
	}
	deliverymen, count, err := s.repo.GetAllDeliverymen(lim, lim*pag)
	if err != nil {
		return nil, 0, err
	}
	return deliverymen, count, nil
}

func (s *DeliverymanService) GetDeliveryman(deliveryId string) (*model.DeliveryManInfoCreate, error) {
	id, err := strconv.Atoi(deliveryId)
	if err != nil {
		return nil, err
	}
	deliveryman, err := s.repo.GetDeliveryman(id)
	if err != nil {
		return nil, err
	}
	return deliveryman, nil
}

func (s *DeliverymanService) CreateDeliveryman(newDeliveryman model.DeliveryManInfoCreate, userId string) (int, error) {
	id, err := strconv.Atoi(userId)
	if err != nil {
		return 0, err
	}
	deliverymanId, err := s.repo.CreateDeliveryman(newDeliveryman, id)
	if err != nil {
		return 0, err
	}
	return deliverymanId, nil
}

func (s *DeliverymanService) ChangeDeliverymanData(newData model.DeliveryManInfoCreate, deliverymanId string) error {
	id, err := strconv.Atoi(deliverymanId)
	if err != nil {
		return err
	}
	err = s.repo.ChangeDeliverymanData(newData, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *DeliverymanService) DeleteDeliveryman(deliverymanId string) error {
	id, err := strconv.Atoi(deliverymanId)
	if err != nil {
		return err
	}
	err = s.repo.DeleteDeliveryman(id)
	if err != nil {
		return err
	}
	return nil
}
