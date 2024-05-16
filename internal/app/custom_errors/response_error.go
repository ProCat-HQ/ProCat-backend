package custom_errors

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/sirupsen/logrus"
)

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error("method: ", c.Request.Method, ", url: ", c.Request.URL, ", statusCode: ", statusCode, ", msg: ", message)

	c.AbortWithStatusJSON(statusCode, model.Response{
		Status:  statusCode,
		Message: message,
		Payload: nil,
	})
}
