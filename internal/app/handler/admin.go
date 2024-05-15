package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
)

func (h *Handler) GetAllDeliveriesToSort(c *gin.Context) {
	count, rows, err := h.services.GetDeliveriesToSort()
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"count": count,
			"rows":  rows,
		},
	})
}

func (h *Handler) ChangeDeliveryData(c *gin.Context) {
	var input struct {
		DeliveryId    int `json:"deliveryId"`
		DeliverymanId int `json:"deliverymanId"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.ChangeDeliveryman(input.DeliveryId, input.DeliverymanId); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
	})
}

func (h *Handler) Cluster(c *gin.Context) {
	payload, err := h.services.Admin.MakeClustering()
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Clustering hasn't done: %s", err))
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"result": payload,
		},
	})
}
