package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetAllDeliveriesToSort(c *gin.Context) {

}

func (h *Handler) ChangeDeliveryData(c *gin.Context) {
}

func (h *Handler) Cluster(c *gin.Context) {
	payload, err := h.services.Admin.MakeClustering()
	if err != nil {
		//custom_errors.NewErrorResponse(c, http.StatusInternalServerError, "Clustering hasn't done")
		logrus.Error("Clustering hasn't done")
		return
	}

	logrus.Info(payload)

	c.JSON(200, gin.H{
		"message": "ok",
		"payload": payload,
	})
}
