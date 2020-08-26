package handler

import (
	"net/http"

	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

type Config struct {
	Store models.Config
}

func (h Config) Get(c echo.Context) error {
	return c.JSON(http.StatusOK, h.Store)
}

func (h Config) Register(app echoswagger.ApiRoot) echoswagger.ApiRoot {
	app.GET("/api/appconfig", h.Get).
		AddResponse(http.StatusOK, "success", models.Config{}, nil)

	return app
}
