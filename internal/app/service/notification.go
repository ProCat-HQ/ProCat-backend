package service

import (
	"errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
)

type NotificationService struct {
	repo repository.Notification
}

func NewNotificationService(repo repository.Notification) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) GetUsersNotification(userId int) ([]model.Notification, error) {
	return s.repo.GetUsersNotification(userId)
}

func (s *NotificationService) CreateNotification(userId int, title, description string) (int, error) {
	return s.repo.CreateNotification(userId, title, description)
}

func (s *NotificationService) ReadAndGetNotification(userId int, notificationId int) (model.Notification, error) {
	userNotificationId, err := s.repo.GetNotificationUserId(notificationId)
	if err != nil {
		return model.Notification{}, err
	}
	if userNotificationId != userId {
		return model.Notification{}, errors.New("can't read someone else's notification")
	}
	return s.repo.ReadAndGetNotification(notificationId)
}

func (s *NotificationService) DeleteNotification(notificationId int) error {
	return s.repo.DeleteNotification(notificationId)
}
