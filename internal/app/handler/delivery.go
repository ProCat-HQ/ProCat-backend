package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/procat-hq/procat-backend/internal/routing"
	"io"
	"net/http"
	"os"
	"strings"
)

func (h *Handler) GetAllDeliveries(c *gin.Context) {

}

func (h *Handler) GetAllDeliveriesForOneDeliveryman(c *gin.Context) {

}

func (h *Handler) ChangeDeliveryStatus(c *gin.Context) {

}

func (h *Handler) CreateRoute(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	requestBody, err := h.services.Delivery.GetDeliveriesForDeliveryman(userData.UserId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	jsn, err := json.Marshal(requestBody)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := http.Post(fmt.Sprintf("https://routing.api.2gis.com/get_dist_matrix?key=%s&version=2.0", os.Getenv("API_KEY_2GIS")),
		"application/json", strings.NewReader(string(jsn)))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := &routing.Response{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result, _ := routing.GetRoute(*response)

	var realResponse []model.LatLon
	for _, val := range result.OptimalRoute {
		realResponse = append(realResponse, requestBody.Points[val])
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"points": realResponse,
		},
	})

	defer resp.Body.Close()
}
