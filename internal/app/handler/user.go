package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/errors"
	"github.com/procat-hq/procat-backend/internal/app/model"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) GetAllUsers(c *gin.Context) {

}

func (h *Handler) GetUser(c *gin.Context) {

}

func (h *Handler) SignIn(c *gin.Context) {

}

func (h *Handler) SignUp(c *gin.Context) {
	var input model.User
	if err := c.BindJSON(&input); err != nil {
		errors.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logrus.Info(fmt.Sprintf("%+v", input))
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
