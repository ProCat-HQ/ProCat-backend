package service

import "github.com/procat-hq/procat-backend/internal/app/repository"

type AdminService struct {
	repo repository.Admin
}

func NewAdminService(repo repository.Admin) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) MakeClustering() error {
	return nil
}
