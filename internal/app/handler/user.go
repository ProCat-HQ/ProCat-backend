package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/custom_errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"net/http"
	"regexp"
	"strings"
)

const (
	phoneRegex = "^\\+?[1-9]?[0-9]{7,14}$"
)

func (h *Handler) GetAllUsers(c *gin.Context) {

}

func (h *Handler) GetUser(c *gin.Context) {

}

type SignInInput struct {
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

func (h *Handler) SignIn(c *gin.Context) {
	var input SignInInput
	if err := c.BindJSON(&input); err != nil {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if matched, _ := regexp.MatchString(phoneRegex, input.PhoneNumber); !matched {
		custom_errors.NewErrorResponse(c, http.StatusBadRequest, "Invalid phoneNumber field")
		return
	}

	token, err := h.services.User.GenerateToken(input.PhoneNumber, input.Password)
	if err != nil {
		custom_errors.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *Handler) SignUp(c *gin.Context) {
	var input model.SignUpInput
	if err := c.BindJSON(&input); err != nil {
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

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) ChangeIIN(c *gin.Context) {

}

func (h *Handler) ChangeFullName(c *gin.Context) {

}

func (h *Handler) ChangePassword(c *gin.Context) {

}

func (h *Handler) ChangePhone(c *gin.Context) {

}

func (h *Handler) ChangeEmail(c *gin.Context) {

}

func (h *Handler) ChangeRole(c *gin.Context) {

}
