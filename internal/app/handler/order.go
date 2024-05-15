package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
	"time"
)

func (h *Handler) GetAllOrders(c *gin.Context) {

}

func (h *Handler) GetOrder(c *gin.Context) {

}

func (h *Handler) CreateOrder(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input model.OrderCreation

	if err := c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	rpStart, err := time.Parse(time.DateTime, input.RentalPeriodStart)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	rpEnd, err := time.Parse(time.DateTime, input.RentalPeriodEnd)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tStart, err := time.Parse(time.DateTime, input.TimeStart)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tEnd, err := time.Parse(time.DateTime, input.TimeEnd)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	orderWithTime := model.OrderCreationWithTime{
		RentalPeriodStart: rpStart,
		RentalPeriodEnd:   rpEnd,
		Address:           input.Address,
		CompanyName:       input.CompanyName,
		DeliveryMethod:    input.DeliveryMethod,
		TimeStart:         tStart,
		TimeEnd:           tEnd,
	}

	orderCheque, err := h.services.Order.CreateOrder(userData.UserId, orderWithTime)

	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: orderCheque,
	})
}

func (h *Handler) CancelOrder(c *gin.Context) {

}

func (h *Handler) ChangeOrderStatus(c *gin.Context) {

}

func (h *Handler) GetPaymentData(c *gin.Context) {

}

func (h *Handler) ChangePaymentStatus(c *gin.Context) {

}
