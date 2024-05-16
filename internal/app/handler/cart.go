package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
	"strconv"
)

func (h *Handler) GetCartItems(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.services.Cart.GetCartItems(userData.UserId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"items": items,
		},
	})
}

type AddToCartStruct struct {
	ItemId int `json:"itemId" binding:"required,min=1"`
	Count  int `json:"count"`
}

func (h *Handler) AddItemsToCart(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input AddToCartStruct
	if err := c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Cart.AddItemsToCart(userData.UserId, input.ItemId, input.Count)
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

func (h *Handler) DeleteItemsFromCart(c *gin.Context) {
	itemIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	count := c.Query("count")
	if count == "" {
		count = "1"
	}

	countInt, err := strconv.Atoi(count)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Cart.DeleteItemFromCart(userData.UserId, itemIdInt, countInt)
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
