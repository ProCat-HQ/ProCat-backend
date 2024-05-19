package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/app/repository"
	"github.com/procat-hq/procat-backend/internal/app/utils"
	"mime/multipart"
	"strconv"
	"strings"
)

type ItemService struct {
	repo repository.Item
}

func NewItemService(repo repository.Item) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) GetAllItems(limit, page, search, categoryId, stock string) (int, []model.PieceOfItem, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return 0, nil, err
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return 0, nil, err
	}

	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		return 0, nil, err
	}

	isInStock, err := strconv.ParseBool(stock)
	if err != nil {
		return 0, nil, err
	}

	offset := pageInt * limitInt

	count, items, err := s.repo.GetAllItems(limitInt, offset, categoryIdInt, isInStock, search)
	if err != nil {
		return 0, nil, err
	}

	return count, items, nil
}

func (s *ItemService) GetItem(itemId string) (model.Item, error) {
	itemIdInt, err := strconv.Atoi(itemId)

	var item model.Item
	if err != nil {
		return item, err
	}

	item, err = s.repo.GetItem(itemIdInt)
	return item, err
}

func (s *ItemService) CreateItem(name, description, price, priceDeposit, categoryId string, files []*multipart.FileHeader) (int, error) {
	priceInt, err := strconv.Atoi(price)
	if err != nil {
		return 0, err
	}
	priceDepositInt, err := strconv.Atoi(priceDeposit)
	if err != nil {
		return 0, err
	}
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		return 0, err
	}

	filenames := make([]string, 0)

	for _, file := range files {
		filenameParts := strings.Split(strings.TrimSpace(file.Filename), ".")
		newFilename := uuid.New().String() + "." + filenameParts[len(filenameParts)-1]

		err := utils.SaveUploadedFileAndCheckExtension(file, "./assets/"+newFilename)
		if err != nil {
			fileError := errors.New("Error while uploading " + file.Filename + ": " + err.Error())
			err1 := utils.RemoveFiles(filenames, "./assets/")
			if err1 != nil {
				return 0, errors.New(fileError.Error() + " and " + err1.Error())
			}
			return 0, fileError
		}
		filenames = append(filenames, newFilename)
	}

	itemId, err := s.repo.CreateItem(name, description, priceInt, priceDepositInt, categoryIdInt)
	if err != nil {
		err1 := utils.RemoveFiles(filenames, "./assets/")
		if err1 != nil {
			return 0, errors.New(err.Error() + " and " + err1.Error())
		}
		return 0, err
	}
	err = s.repo.SaveFilenames(itemId, filenames)
	if err != nil {
		err1 := utils.RemoveFiles(filenames, "./assets/")
		if err1 != nil {
			return 0, errors.New(err.Error() + " and " + err1.Error())
		}
		return 0, err
	}
	return itemId, nil
}

func (s *ItemService) ChangeItem(itemId int, name, description, price, priceDeposit, categoryId *string) error {
	if price != nil {
		_, err := strconv.Atoi(*price)
		if err != nil {
			return err
		}
	}

	if priceDeposit != nil {
		_, err := strconv.Atoi(*priceDeposit)
		if err != nil {
			return err
		}
	}

	if categoryId != nil {
		_, err := strconv.Atoi(*categoryId)
		if err != nil {
			return err
		}
	}

	return s.repo.ChangeItem(itemId, name, description, price, priceDeposit, categoryId)
}

func (s *ItemService) DeleteItem(itemId int) error {
	return s.repo.DeleteItem(itemId)
}

func (s *ItemService) ChangeStockOfItem(itemId, storeId, inStockNumber int) error {
	if inStockNumber < 0 {
		return errors.New("inStockNumber must be positive")
	}
	return s.repo.ChangeStockOfItem(itemId, storeId, inStockNumber)
}
