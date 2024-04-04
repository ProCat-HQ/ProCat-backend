package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) GetAllDeliveries(c *gin.Context) {

}

func (h *Handler) GetDelivery(c *gin.Context) {

}

func (h *Handler) ChangeDeliveryStatus(c *gin.Context) {

}

func (h *Handler) CreateRoute(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	logrus.Info(gin.H{
		"userId=": userData.UserId,
	})

	requestBody, err := h.services.Delivery.GetDeliveriesForDeliveryman(userData.UserId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//json, err := json2.Marshal(requestBody)
	//if err != nil {
	//	custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}

	logrus.Info("SUCCESS")

	//resp, err := http.Post("https://routing.api.2gis.com/get_dist_matrix?key=810e358b-1439-4919-9eab-4618b85be168&version=2.0",
	//	"application/json", strings.NewReader(string(json)))
	//if err != nil {
	//	custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//
	//if err != nil {
	//	custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}

	c.JSON(http.StatusOK, requestBody)

	//defer resp.Body.Close()
}
