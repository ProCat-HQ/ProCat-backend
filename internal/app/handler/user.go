package handler

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	phoneRegex = "^\\+?[1-9]?[0-9]{7,14}$"
	emailRegex = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
)

func (h *Handler) GetAllUsers(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "20"
	}
	page := c.Query("page")
	if page == "" {
		page = "0"
	}
	role := c.Query("role")
	isConfirmed := c.Query("isConfirmed")

	count, users, err := h.services.User.GetAllUsers(limit, page, role, isConfirmed)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"count": count,
			"rows":  users,
		},
	})
}

func (h *Handler) GetUser(c *gin.Context) {
	paramUserId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "userId param is not a number: "+err.Error())
		return
	}

	user, err := h.services.User.GetUserById(paramUserId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: user,
	})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	paramUserId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "userId param is not a number: "+err.Error())
		return
	}

	selfId, ok := c.Get("userId")
	if !ok {
		custom_errors.NewErrorResponse(c, http.StatusUnauthorized, "userId field not found in context")
		return
	}
	if selfId == paramUserId {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "please don't delete yourself")
		return
	}

	err = h.services.User.DeleteUserById(paramUserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			custom_errors.NewErrorResponse(c, http.StatusBadRequest, "no user with such id")
			return
		}
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

type SignInInput struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Fingerprint string `json:"fingerprint" binding:"required"`
}

func (h *Handler) SignIn(c *gin.Context) {
	var input SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if matched, _ := regexp.MatchString(phoneRegex, input.PhoneNumber); !matched {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Invalid phoneNumber field")
		return
	}

	user, err := h.services.User.GetUserByCredentials(input.PhoneNumber, input.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			custom_errors.NewErrorResponse(c, http.StatusUnauthorized, "bad credentials: wrong phone number or password")
			return
		}
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, refreshToken, err := h.services.User.GenerateTokens(user, input.Fingerprint)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	})
}

func (h *Handler) SignUp(c *gin.Context) {
	var input model.SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if len(strings.Split(input.FullName, " ")) < 2 {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Invalid fullName field")
		return
	}

	if matched, _ := regexp.MatchString(phoneRegex, input.PhoneNumber); !matched {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Invalid phoneNumber field")
		return
	}

	id, err := h.services.User.CreateUser(input)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"id": id,
		},
	})
}

func (h *Handler) Logout(c *gin.Context) {
	var refreshToken struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	if err := c.ShouldBindJSON(&refreshToken); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	refreshClaims, err := h.services.User.ParseRefreshToken(refreshToken.RefreshToken)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	status, err := h.services.User.LogoutUser(refreshToken.RefreshToken, refreshClaims.UserId)
	if err != nil {
		custom_errors.NewErrorResponse(c, status, err.Error())
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

type RefreshTokenStruct struct {
	RefreshToken string `json:"refreshToken"`
	Fingerprint  string `json:"fingerprint"`
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var refreshTokenStruct RefreshTokenStruct
	if err := c.ShouldBindJSON(&refreshTokenStruct); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	claims, err := h.services.User.ParseRefreshToken(refreshTokenStruct.RefreshToken)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	accessToken, refreshToken, err := h.services.User.RegenerateTokens(claims.UserId, refreshTokenStruct.RefreshToken, refreshTokenStruct.Fingerprint)
	if err != nil {
		if strings.Contains(err.Error(), "suspicious activity") {
			custom_errors.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	})
}

func (h *Handler) ChangeIIN(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input struct {
		Password             string `json:"password" binding:"required"`
		IdentificationNumber string `json:"identificationNumber" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	passwordMatch, err := h.services.User.CheckPassword(input.Password, userData.UserId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !passwordMatch {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "password does not match")
		return
	}

	if err = h.services.User.ChangeIdentificationNumber(userData.UserId, input.IdentificationNumber); err != nil {
		if err.Error() == "wrong iinBin format" {
			custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) ChangeFullName(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input struct {
		Password string `json:"password" binding:"required"`
		FullName string `json:"fullName" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if len(strings.Split(input.FullName, " ")) < 2 {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Invalid fullName field")
		return
	}

	passwordMatch, err := h.services.User.CheckPassword(input.Password, userData.UserId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !passwordMatch {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "password does not match")
		return
	}

	if err = h.services.User.ChangeFullName(userData.UserId, input.FullName); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	passwordMatch, err := h.services.User.CheckPassword(input.OldPassword, userData.UserId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !passwordMatch {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "password does not match")
		return
	}
	if err = h.services.User.ChangePassword(userData.UserId, input.NewPassword); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) ChangePhone(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input struct {
		Password    string `json:"password" binding:"required"`
		PhoneNumber string `json:"phoneNumber" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if matched, _ := regexp.MatchString(phoneRegex, input.PhoneNumber); !matched {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Invalid phoneNumber field")
		return
	}

	passwordMatch, err := h.services.User.CheckPassword(input.Password, userData.UserId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !passwordMatch {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "password does not match")
		return
	}

	if err = h.services.User.ChangePhoneNumber(userData.UserId, input.PhoneNumber, input.Password); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) ChangeEmail(c *gin.Context) {
	userData, err := h.GetUserContext(c)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input struct {
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if matched, _ := regexp.MatchString(emailRegex, input.Email); !matched {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Invalid email format field")
		return
	}

	passwordMatch, err := h.services.User.CheckPassword(input.Password, userData.UserId)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !passwordMatch {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "password does not match")
		return
	}

	if err = h.services.User.ChangeEmail(userData.UserId, input.Email); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}

func (h *Handler) ChangeRole(c *gin.Context) {
	paramUserId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "userId param is not a number: "+err.Error())
		return
	}

	var input struct {
		Role string `json:"role" binding:"required"`
	}

	if err = c.ShouldBindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.services.User.ChangeUserRole(paramUserId, input.Role); err != nil {
		if err.Error() == "unknown role" {
			custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "ok",
		Payload: nil,
	})
}
