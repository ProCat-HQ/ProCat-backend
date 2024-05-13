package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
)

func (h *Handler) GetAllDeliverymen(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	page := c.Query("page")
	if page == "" {
		page = "0"
	}

	deliverymen, count, err := h.services.Deliveryman.GetAllDeliverymen(limit, page)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"count": count,
			"rows":  deliverymen,
		},
	})
}

func (h *Handler) GetDeliveryman(c *gin.Context) {
	deliverymanId := c.Param("id")

	payload, err := h.services.Deliveryman.GetDeliveryman(deliverymanId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if payload == nil {
		c.JSON(http.StatusOK, model.Response{
			Status:  http.StatusOK,
			Message: "ok",
			Payload: nil,
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: payload,
	})
}

func (h *Handler) CreateDeliveryman(c *gin.Context) {
	var deliveryman model.DeliveryManInfoCreate
	if err := c.ShouldBindJSON(&deliveryman); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userId := c.Param("id")

	newId, err := h.services.CreateDeliveryman(deliveryman, userId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"id": newId,
		},
	})
}

func (h *Handler) ChangeDeliverymanData(c *gin.Context) {
	var deliveryman model.DeliveryManInfoCreate
	if err := c.ShouldBindJSON(&deliveryman); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	deliverymanId := c.Param("id")

	err := h.services.ChangeDeliverymanData(deliveryman, deliverymanId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
	})
}

func (h *Handler) DeleteDeliveryman(c *gin.Context) {
	deliverymanId := c.Param("id")

	err := h.services.DeleteDeliveryman(deliverymanId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
	})
}
