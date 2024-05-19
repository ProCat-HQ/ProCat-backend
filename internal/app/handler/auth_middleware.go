package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
	"strconv"
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

	userData, err := h.services.User.ParseAccessToken(headerParts[1])
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set("userId", userData.UserId)
	c.Set("userRole", userData.UserRole)
	c.Next()
}

func (h *Handler) MustBelongsToUser(c *gin.Context) {
	paramUserId, err := strconv.Atoi(c.Param("id")) // string
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "userId param is not a number: "+err.Error())
		return
	}

	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	rolePriority, err := getRolePriority(userData.UserRole)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if userData.UserId != paramUserId && rolePriority < 4 {
		custom_errors.NewErrorResponse(c, http.StatusForbidden, "Forbidden")
		return
	}

	c.Next()
}

func (h *Handler) GetUserContext(c *gin.Context) (*model.AccessTokenClaimsExtension, error) {
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

	userContext := &model.AccessTokenClaimsExtension{
		UserId:   userId,
		UserRole: userRole,
	}
	return userContext, nil
}

func getRolePriority(role string) (int, error) {
	switch role {
	case model.UserRole:
		return 1, nil
	case model.DeliverymanRole:
		return 2, nil
	case model.ModeratorRole:
		return 3, nil
	case model.AdminRole:
		return 4, nil
	default:
		return 0, errors.New("unknown role")
	}
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
			custom_errors.NewErrorResponse(c, http.StatusBadRequest, "User role is not a string type")
			return
		}

		userRolePriority, err := getRolePriority(userRole)
		if err != nil {
			custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		roleToCheckPriority, err := getRolePriority(roleToCheck)
		if err != nil {
			custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		if userRolePriority < roleToCheckPriority {
			custom_errors.NewErrorResponse(c, http.StatusForbidden, "Forbidden")
			return
		}

		c.Next()
	}
	return gin.HandlerFunc(fn)
}
