package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
	"strconv"
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

	search := c.Query("search")

	categoryId := c.Query("categoryId")
	if categoryId == "" {
		categoryId = "0" // if == 0 => category won't be included to query
	}

	stock := c.Query("stock")
	if stock == "" {
		stock = "false"
	}

	count, items, err := h.services.Item.GetAllItems(limit, page, search, categoryId, stock)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"count": count,
			"rows":  items,
		},
	})
}

func (h *Handler) GetItem(c *gin.Context) {
	itemId := c.Param("id")
	item, err := h.services.Item.GetItem(itemId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: item,
	})
}

func (h *Handler) CreateItem(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Invalid input form")
		return
	}

	name, ok := form.Value["name"]
	if !ok || name[0] == "" {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "No 'name' field included")
		return
	}
	description := ""
	descriptionList, ok := form.Value["description"]
	if ok {
		description = descriptionList[0]
	}
	price, ok := form.Value["price"]
	if !ok || price[0] == "" {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "No 'price' field included")
		return
	}
	priceDeposit, ok := form.Value["priceDeposit"]
	if !ok || priceDeposit[0] == "" {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "No 'priceDeposit' field included")
		return
	}
	categoryId, ok := form.Value["categoryId"]
	if !ok || categoryId[0] == "" {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "No 'categoryId' field included")
		return
	}

	files, ok := form.File["images"]
	if !ok {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "No files included")
		return
	}
	id, err := h.services.Item.CreateItem(name[0], description, price[0], priceDeposit[0], categoryId[0], files)
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

func (h *Handler) ChangeItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Invalid item id")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Invalid input form")
		return
	}

	var name, description, price, priceDeposit, categoryId *string
	name, description, price, priceDeposit, categoryId = nil, nil, nil, nil, nil

	n, ok := form.Value["name"]
	if ok {
		name = &n[0]
	}
	d, ok := form.Value["description"]
	if ok {
		description = &d[0]
	}
	p, ok := form.Value["price"]
	if ok {
		price = &p[0]
	}
	pd, ok := form.Value["priceDeposit"]
	if ok {
		priceDeposit = &pd[0]
	}
	ca, ok := form.Value["categoryId"]
	if ok {
		categoryId = &ca[0]
	}

	err = h.services.Item.ChangeItem(itemId, name, description, price, priceDeposit, categoryId)
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

func (h *Handler) DeleteItem(c *gin.Context) {
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Invalid item id")
		return
	}

	err = h.services.Item.DeleteItem(itemId)
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

func (h *Handler) AddInfo(c *gin.Context) {
	itemIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input model.ItemInfoCreation
	if err = c.ShouldBind(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Item.AddInfos(itemIdInt, input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) DeleteInfo(c *gin.Context) {
	itemIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ids := c.QueryArray("id")
	if ids == nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "no ids were passed")
		return
	}
	idsInt := make([]int, 0, len(ids))
	for _, id := range ids {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			custom_errors.NewErrorResponse(c, http.StatusBadRequest, "id query param is not a number")
			return
		}
		idsInt = append(idsInt, idInt)
	}

	if err = h.services.Item.DeleteInfos(itemIdInt, idsInt); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) ChangeInfo(c *gin.Context) {
	itemIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input model.ItemInfoChange
	if err = c.ShouldBind(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Item.ChangeInfos(itemIdInt, input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) AddImages(c *gin.Context) {
	itemIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Invalid input form")
		return
	}

	files, ok := form.File["images"]
	if !ok {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "No files included")
		return
	}

	if err = h.services.Item.AddImages(itemIdInt, files); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) DeleteImages(c *gin.Context) {
	itemIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	ids := c.QueryArray("id")
	if ids == nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "no ids were passed")
		return
	}
	idsInt := make([]int, 0, len(ids))
	for _, id := range ids {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			custom_errors.NewErrorResponse(c, http.StatusBadRequest, "id query param is not a number")
			return
		}
		idsInt = append(idsInt, idInt)
	}

	if err = h.services.Item.DeleteImages(itemIdInt, idsInt); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) ChangeStock(c *gin.Context) {
	itemIdInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input model.ChangeStock
	if err = c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.Item.ChangeStockOfItem(itemIdInt, input.StoreId, input.InStockNumber); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}
