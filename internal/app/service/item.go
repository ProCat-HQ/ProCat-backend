package service

import (
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"strconv"
)

type ItemService struct {
	repo repository.Item
}

func NewItemService(repo repository.Item) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) GetAllItems(limit, page, categoryId, stock string) ([]model.PieceOfItem, error) {
	lim, err := strconv.Atoi(limit)
	if err != nil {
		return nil, err
	}
	pag, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}

	catId, err := strconv.Atoi(categoryId)
	if err != nil {
		return nil, err
	}

	isInStock, err := strconv.ParseBool(stock)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.GetAllItems(lim, pag*lim, catId, isInStock)
	if err != nil {
		return nil, err
	}

	return items, nil
}
