package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

const (
	swaggerUI = "https://petstore.swagger.io/?url="
)

func (h *Handler) HandleSwagger(c *gin.Context) {
	url := swaggerUI + "http://79.137.205.181:" + os.Getenv("BIND_ADDR") + "/swagger/docs/" + "api.json"
	c.Redirect(http.StatusTemporaryRedirect, url)
}
