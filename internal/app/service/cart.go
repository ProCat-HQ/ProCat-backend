package service

import (
	"errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
)

type CartService struct {
	repo repository.Cart
}

func NewCartService(repo repository.Cart) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) AddItemsToCart(userId, itemId, count int) error {
	cartId, err := s.repo.GetUsersCartId(userId)
	if err != nil {
		return err
	}

	if count < 0 {
		return errors.New("count must be greater than zero")
	}

	err = s.repo.AddItemToCart(cartId, itemId, count)

	return err
}

func (s *CartService) DeleteItemFromCart(userId, itemId int) error {
	cartId, err := s.repo.GetUsersCartId(userId)
	if err != nil {
		return err
	}

	if itemId <= 0 {
		return errors.New("itemId must be positive")
	}

	err = s.repo.DeleteItemFromCart(cartId, itemId)
	return err
}

func (s *CartService) GetCartItems(userId int) ([]model.CartItem, error) {
	cartId, err := s.repo.GetUsersCartId(userId)
	if err != nil {
		return nil, err
	}

	return s.repo.GetCartItems(cartId)
}
