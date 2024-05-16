package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) GetAllOrders(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	page := c.Query("page")
	if page == "" {
		page = "0"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	rolePriority, err := getRolePriority(userData.UserRole)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId := c.Query("userId")
	var userIdInt int
	if userId == "" {
		if rolePriority >= 4 {
			userIdInt = 0 // => get ALL orders
		} else {
			userIdInt = userData.UserId
		}
	} else {
		userIdInt, err = strconv.Atoi(userId)
		if err != nil {
			custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		if userIdInt <= 0 {
			custom_errors.NewErrorResponse(c, http.StatusBadRequest, "userId must be greater than zero")
			return
		}
		if userData.UserId != userIdInt && rolePriority < 4 {
			custom_errors.NewErrorResponse(c, http.StatusForbidden, "Forbidden")
			return
		}
	}

	statuses := c.QueryArray("status")
	count, rows, err := h.services.Order.GetAllOrders(limitInt, pageInt, userIdInt, statuses)
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
