package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"net/http"
	"strings"
)

const (
	authHeader = "Authorization"
)

func (h *Handler) UserIdentify(c *gin.Context) {
	header := c.GetHeader(authHeader)
	if header == "" {
		custom_errors.NewErrorResponse(c, http.StatusUnauthorized, "Empty Authorization header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		custom_errors.NewErrorResponse(c, http.StatusUnauthorized, "Invalid Authorization header")
		return
	}

	if len(headerParts[1]) == 0 {
		custom_errors.NewErrorResponse(c, http.StatusUnauthorized, "Empty token")
		return
	}
}
