package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
)

func (h *Handler) GetAllItems(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	page := c.Query("page")
	if page == "" {
		page = "0"
	}
	categoryId := c.Query("categoryId")
	if categoryId == "" {
		categoryId = "0"
	}
	stock := c.Query("stock")
	if stock == "" {
		stock = "false"
	}

	items, err := h.services.Item.GetAllItems(limit, page, categoryId, stock)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var itemsToRes []model.PieceOfItemToRes

	for _, item := range items {
		var desc *string
		if !item.Description.Valid {
			desc = nil
		} else {
			desc = &item.Description.String
		}

		var cat *int
		if !item.CategoryId.Valid {
			cat = nil
		} else {
			*cat = int(item.CategoryId.Int32)
		}
		i := model.PieceOfItemToRes{
			Id:          item.Id,
			Name:        item.Name,
			Description: desc,
			Price:       item.Price,
			IsInStock:   false, // TODO
			Images:      nil,
			CategoryId:  cat,
		}

		itemsToRes = append(itemsToRes, i)
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"count": 120, // TODO
			"rows":  itemsToRes,
		},
	})

}

func (h *Handler) GetItem(c *gin.Context) {

}

func (h *Handler) CreateItem(c *gin.Context) {
	// TODO: не работает, переделать
	form, err := c.MultipartForm()
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Invalid input form")
		return
	}
	files, ok := form.File["upload[]"]
	if !ok {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "No files included")
		return
	}
	for _, file := range files {
		err := c.SaveUploadedFile(file, "./static")
		if err != nil {
			custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Error while uploading "+file.Filename)
			return
		}
	}
	c.String(http.StatusOK, fmt.Sprintf("OK"))
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

func (h *Handler) DeleteImages(c *gin.Context) {

}

func (h *Handler) ChangeStock(c *gin.Context) {

}
