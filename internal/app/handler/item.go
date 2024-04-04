package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"net/http"
)

func (h *Handler) GetAllItems(c *gin.Context) {
	id, ok := c.Get("userId")
	if !ok {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "user id not found")
		return
	}
	role, ok := c.Get("userRole")
	if !ok {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "user role not found")
		return
	}

	c.JSON(200, gin.H{
		"id":   id,
		"role": role,
	})
}

func (h *Handler) GetItem(c *gin.Context) {

}

func (h *Handler) CreateItem(c *gin.Context) {

}

func (h *Handler) ChangeItem(c *gin.Context) {

}

func (h *Handler) DeleteItem(c *gin.Context) {

}

func (h *Handler) AddInfo(c *gin.Context) {

}

func (h *Handler) DeleteInfo(c *gin.Context) {

}

func (h *Handler) ChangeInfo(c *gin.Context) {

}

func (h *Handler) AddImages(c *gin.Context) {

}

func (h *Handler) ChangeImages(c *gin.Context) {

}

func (h *Handler) DeleteImages(c *gin.Context) {

}

func (h *Handler) CreateStock(c *gin.Context) {

}

func (h *Handler) ChangeStock(c *gin.Context) {

}
