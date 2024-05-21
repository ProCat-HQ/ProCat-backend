package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
	"strconv"
)

func (h *Handler) GetAllNotifications(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	notifications, err := h.services.Notification.GetUsersNotification(userData.UserId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"notifications": notifications,
		},
	})
}

func (h *Handler) SendNotification(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input model.NotificationCreate
	if err = c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Notification.CreateNotification(userId, input.Title, input.Description)
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

func (h *Handler) ViewNotification(c *gin.Context) {
	notificationId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	notification, err := h.services.Notification.ReadAndGetNotification(userData.UserId, notificationId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: notification,
	})
}

func (h *Handler) DeleteNotification(c *gin.Context) {
	notificationId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Notification.DeleteNotification(notificationId); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}
