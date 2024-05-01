package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
)

func (h *Handler) GetAllDeliveriesToSort(c *gin.Context) {

}

func (h *Handler) ChangeDeliveryData(c *gin.Context) {
}

func (h *Handler) Cluster(c *gin.Context) {
	payload, err := h.services.Admin.MakeClustering()
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("Clustering hasn't done: %s", err))
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"result": payload,
		},
	})
}
