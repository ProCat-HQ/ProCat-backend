package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) GetAllStores(c *gin.Context) {
	stores, err := h.services.Store.GetAllStores()
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: stores,
	})
}

func (h *Handler) CreateStore(c *gin.Context) {
	var storeData model.StoreCreation

	if err := c.ShouldBindJSON(&storeData); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	timeStart, err := time.Parse(time.TimeOnly, storeData.WorkingHoursStart)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	timeEnd, err := time.Parse(time.TimeOnly, storeData.WorkingHoursEnd)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	store := model.Store{
		Name:              storeData.Name,
		Address:           storeData.Address,
		WorkingHoursStart: timeStart,
		WorkingHoursEnd:   timeEnd,
	}

	id, err := h.services.Store.CreateStore(store)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"id": id,
		},
	})

}

func (h *Handler) ChangeStore(c *gin.Context) {
	storeIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var storeData model.StoreChange

	if err = c.ShouldBindJSON(&storeData); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Store.ChangeStore(storeIdInt, storeData); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) DeleteStore(c *gin.Context) {
	storeIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Store.DeleteStore(storeIdInt); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}
