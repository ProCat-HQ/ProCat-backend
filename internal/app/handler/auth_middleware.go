package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
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

	userData, err := h.services.User.ParseToken(headerParts[1])
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set("userId", userData.UserId)
	c.Set("userRole", userData.UserRole)
	c.Next()
}

func (h *Handler) GetUserContext(c *gin.Context) (*model.TokenClaimsExtension, error) {
	id, ok := c.Get("userId")
	if !ok {
		return nil, errors.New("userId field not found in context")
	}

	role, ok := c.Get("userRole")
	if !ok {
		return nil, errors.New("userRole field not found in context")
	}

	userId, ok := id.(int)
	if !ok {
		return nil, errors.New("userId is not a numeric (integer) type")
	}

	userRole, ok := role.(string)
	if !ok {
		return nil, errors.New("userRole is not a string type")
	}

	userContext := &model.TokenClaimsExtension{
		UserId:   userId,
		UserRole: userRole,
	}
	return userContext, nil
}

func (h *Handler) CheckRole(roleToCheck string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		role, ok := c.Get("userRole")
		if !ok {
			custom_errors.NewErrorResponse(c, http.StatusUnauthorized, "Can't get userRole data")
			return
		}

		userRole, ok := role.(string)
		if !ok {
			custom_errors.NewErrorResponse(c, http.StatusUnauthorized, "User role is not a string type")
			return
		}

		if userRole != roleToCheck {
			custom_errors.NewErrorResponse(c, http.StatusForbidden, "Forbidden")
		}

		c.Next()
	}
	return fn
}
