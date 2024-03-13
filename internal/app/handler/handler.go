package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/procat-hq/procat-backend/internal/app/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	user := router.Group("/user")
	{
		user.POST("/sign-up", h.signUp)
		user.POST("/sign-in", h.signIn)
	}

	api := router.Group("/api")
	{
		api.POST("/route", h.getRoute)
	}
	return router
}
