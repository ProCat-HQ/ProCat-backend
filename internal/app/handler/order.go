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
	orderId, err := strconv.Atoi(c.Param("id"))
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

	order, err := h.services.Order.GetOrder(orderId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userData.UserId != order.UserId && rolePriority < 2 {
		custom_errors.NewErrorResponse(c, http.StatusForbidden, "Forbidden")
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: order,
	})
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
	orderId, err := strconv.Atoi(c.Param("id"))
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

	order, err := h.services.Order.GetOrder(orderId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userData.UserId != order.UserId && rolePriority < 4 {
		custom_errors.NewErrorResponse(c, http.StatusForbidden, "Forbidden")
		return
	}

	payments, err := h.services.Order.GetPaymentsForOrder(orderId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	paidSum := 0
	for _, payment := range payments {
		paidSum += payment.Paid
	}
	var newStatus string
	if paidSum > 0 {
		newStatus = model.AwaitingMoneyBack
	} else {
		newStatus = model.AwaitingRejection
	}

	err = h.services.Order.ChangeOrderStatus(orderId, newStatus)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) ChangeOrderStatus(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Order.ChangeOrderStatus(orderId, req.Status)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) GetPaymentData(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Param("id"))
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

	order, err := h.services.Order.GetOrder(orderId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userData.UserId != order.UserId && rolePriority < 4 {
		custom_errors.NewErrorResponse(c, http.StatusForbidden, "Forbidden")
		return
	}

	payments, err := h.services.Order.GetPaymentsForOrder(orderId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"payments": payments,
		},
	})
}

func (h *Handler) ChangePaymentStatus(c *gin.Context) {
	paymentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var req struct {
		Paid   int    `json:"paid" binding:"required"`
		Method string `json:"method" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Paid <= 0 {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "'paid' field must be greater than zero")
		return
	}

	err = h.services.Order.ChangePaymentStatus(paymentId, req.Paid, req.Method)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) ExtendOrder(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input struct {
		RentalPeriodEnd string `json:"rentalPeriodEnd" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	rentalPeriodEnd, err := time.Parse(time.DateTime, input.RentalPeriodEnd)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.services.Order.GetOrder(orderId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userData.UserId != order.UserId {
		custom_errors.NewErrorResponse(c, http.StatusForbidden, "Forbidden")
		return
	}

	if order.RentalPeriodEnd.After(rentalPeriodEnd) {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "new rental end time must be greater than old")
		return
	}

	if err = h.services.Order.ExtendOrder(orderId, rentalPeriodEnd); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) ConfirmOrderExtension(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.services.Order.GetOrder(orderId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if order.Status != model.ExtensionRequest {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "can't extend order without user request")
		return
	}

	if err = h.services.Order.ConfirmOrderExtension(order); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) ReturnOrder(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input struct {
		Problem           bool   `json:"problem" binding:"required"`
		DeliveryMethod    string `json:"deliveryMethod" binding:"required"`
		DeliveryTimeStart string `json:"deliveryTimeStart" binding:"required"`
		DeliveryTimeEnd   string `json:"deliveryTimeEnd" binding:"required"`
	}

	if err = c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.services.Order.GetOrder(orderId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userData.UserId != order.UserId {
		custom_errors.NewErrorResponse(c, http.StatusForbidden, "Forbidden")
		return
	}

	if err = h.services.Order.ReturnOrder(orderId, input.Problem, input.DeliveryMethod,
		input.DeliveryTimeStart, input.DeliveryTimeEnd); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) NeedRepairForOrder(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input struct {
		Price int `json:"price" binding:"required"`
	}

	if err = c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.services.Order.GetOrder(orderId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if order.Status != model.Returned {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "can't need repair if order wasn't returned")
		return
	}

	if err = h.services.Order.NeedRepairForOrder(orderId, input.Price); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}
