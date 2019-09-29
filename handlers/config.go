package handlers

import (
	"net/http"

	"github.com/eiladin/go-simple-startpage/config"
	"github.com/gin-gonic/gin"
)

func GetConfigHandler(c *gin.Context) {
	c.JSON(http.StatusOK, config.GetConfig())
}
