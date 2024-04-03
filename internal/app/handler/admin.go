package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"net/http"
)

func (h *Handler) GetAllDeliveriesToSort(c *gin.Context) {

}

func (h *Handler) ChangeDeliveryData(c *gin.Context) {

}

func (h *Handler) Cluster(c *gin.Context) {
	err := h.services.Admin.MakeClustering()

	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, "Clustering hasn't done")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
