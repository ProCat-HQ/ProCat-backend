package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
)

func (h *Handler) getRoute(c *gin.Context) {
	var input model.RouteList

	if err := c.BindJSON(&input); err != nil {
		errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//h.
}
