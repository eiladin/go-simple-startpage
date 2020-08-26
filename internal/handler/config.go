package handler

import (
	"net/http"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

// Config struct
type Config struct {
	Store config.Config
}

// Get /api/appconfig
func (h Config) Get(c echo.Context) error {
	return c.JSON(http.StatusOK, h.Store)
}

// Register handler
func (h Config) Register(app echoswagger.ApiRoot) echoswagger.ApiRoot {
	app.GET("/api/appconfig", h.Get).
		AddResponse(http.StatusOK, "success", config.Config{}, nil)

	return app
}
