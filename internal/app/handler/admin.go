package handler

import "github.com/gin-gonic/gin"

func (h *Handler) GetAllDeliveriesToSort(c *gin.Context) {

}

func (h *Handler) ChangeDeliveryData(c *gin.Context) {

}

func (h *Handler) Cluster(c *gin.Context) {
	err := h.services.Admin.MakeClustering()

	if err != nil {

	}
}
