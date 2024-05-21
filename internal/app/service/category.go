package service

import (
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
)

type CategoryService struct {
	repo repository.Category
}

func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(categoryParentId int, name string) (int, error) {
	return s.repo.CreateCategory(categoryParentId, name)
}

func (s *CategoryService) ChangeCategory(categoryId int, name string) error {
	return s.repo.ChangeCategory(categoryId, name)
}

func (s *CategoryService) GetCategoriesForParent(categoryParentId int) ([]model.Category, error) {
	return s.repo.GetCategoriesForParent(categoryParentId)
}

func (s *CategoryService) DeleteCategory(categoryId int) error {
	return s.repo.DeleteCategory(categoryId)
}

func (s *CategoryService) GetCategoryRoute(categoryId int) ([]model.Category, error) {
	return s.repo.GetCategoryRoute(categoryId)
}
