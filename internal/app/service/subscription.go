package service

import (
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"strconv"
)

type SubscriptionService struct {
	repo repository.Subscription
}

func NewSubscriptionService(repo repository.Subscription) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) GetUserSubscriptions(userId int, limit, page string) (int, []model.Subscription, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return 0, nil, err
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return 0, nil, err
	}

	offset := limitInt * pageInt
	return s.repo.GetUserSubscriptions(userId, limitInt, offset)
}

func (s *SubscriptionService) CreateSubscription(userId, itemId int) error {
	return s.repo.CreateSubscription(userId, itemId)
}

func (s *SubscriptionService) DeleteSubscription(userId, subId int) error {
	return s.repo.DeleteSubscription(userId, subId)
}
