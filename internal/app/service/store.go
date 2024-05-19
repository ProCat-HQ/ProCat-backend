package service

import (
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"github.com/procat-hq/procat-backend/internal/twogis"
	"time"
)

type StoreService struct {
	repo repository.Store
}

func NewStoreService(repo repository.Store) *StoreService {
	return &StoreService{repo: repo}
}

func (s *StoreService) CreateStore(store model.Store) (int, error) {
	lat, lon, err := twogis.GetLatLon(store.Address)
	if err != nil {
		return -1, err
	}

	store.Latitude = lat
	store.Longitude = lon

	return s.repo.CreateStore(store)
}

func (s *StoreService) GetAllStores() ([]model.StoreFromDB, error) {
	return s.repo.GetAllStores()
}

func (s *StoreService) ChangeStore(storeId int, store model.StoreChange) error {
	storeDataToDB := model.StoreChangeDB{
		Name:    store.Name,
		Address: store.Address,
	}

	if store.WorkingHoursStart != "" {
		timeStart, err := time.Parse(time.TimeOnly, store.WorkingHoursStart)
		if err != nil {
			return err
		}
		storeDataToDB.WorkingHoursStart = &timeStart
	}

	if store.WorkingHoursEnd != "" {
		timeEnd, err := time.Parse(time.TimeOnly, store.WorkingHoursEnd)
		if err != nil {
			return err
		}
		storeDataToDB.WorkingHoursStart = &timeEnd
	}

	if store.Address != "" {
		lat, lon, err := twogis.GetLatLon(store.Address)
		if err != nil {
			return err
		}
		storeDataToDB.Latitude, storeDataToDB.Longitude = lat, lon
	}

	return s.repo.ChangeStore(storeId, storeDataToDB)
}

func (s *StoreService) DeleteStore(storeId int) error {
	return s.repo.DeleteStore(storeId)
}
