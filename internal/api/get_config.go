package api

import (
	"net/http"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

func (h handler) getConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, h.Config)
}

func (h handler) addGetConfigRoute(app echoswagger.ApiRoot) echoswagger.ApiRoot {
	app.GET("/api/appconfig", h.getConfig).
		AddResponse(http.StatusOK, "success", models.Config{}, nil)

	return app
}
