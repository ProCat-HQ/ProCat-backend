package service

import (
	"errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"github.com/procat-hq/procat-backend/internal/kzgov"
	"github.com/procat-hq/procat-backend/internal/twogis"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(userId int, order model.OrderCreationWithTime) (model.OrderCheque, error) {
	if order.RentalPeriodStart.After(order.RentalPeriodEnd) {
		return model.OrderCheque{}, errors.New("rentalPeriod wrong period order")
	}
	if order.TimeStart.After(order.TimeEnd) {
		return model.OrderCheque{}, errors.New("time wrong period order")
	}

	duration := order.RentalPeriodEnd.Sub(order.RentalPeriodStart).Hours()
	if duration <= 0 {
		return model.OrderCheque{}, errors.New("rental period of order less or equals to zero")
	}

	rentPeriodDays := int(duration/24) + 1

	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return model.OrderCheque{}, err
	}

	if user.IdentificationNumber == "" {
		return model.OrderCheque{}, errors.New("user identification number is empty")
	}
	arrearResponse, err := kzgov.GetArrear(user.IdentificationNumber)
	if err != nil {
		return model.OrderCheque{}, err
	}

	if !kzgov.CompareNames(arrearResponse.NameKk, arrearResponse.NameRu, user.FullName) {
		return model.OrderCheque{}, errors.New("fullname from kz.gov service doesn't match with user's fullname")
	}

	lat, lon, err := twogis.GetLatLon(order.Address)
	if err != nil {
		return model.OrderCheque{}, err
	}

	defaultStatus := "awaitingPayment"

	deposit := arrearResponse.TotalArrear > 0

	orderCheque, err := s.repo.CreateOrder(defaultStatus, deposit, order.RentalPeriodStart, order.RentalPeriodEnd,
		order.Address, lat, lon, order.CompanyName, userId, order.DeliveryMethod,
		order.TimeStart, order.TimeEnd, rentPeriodDays)

	return orderCheque, err
}

func (s *OrderService) GetAllOrders(limit, page, userId int, statuses []string) (int, []model.Order, error) {
	offset := limit * page
	count, rows, err := s.repo.GetAllOrders(limit, offset, userId, statuses)
	if err != nil {
		return 0, nil, err
	}
	return count, rows, nil
}
