package handler

import (
	"net/http"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/labstack/echo/v4"
)

// Config struct
type Config struct {
	Store config.Config
}

// Get /api/appconfig
func (h Config) Get(c echo.Context) error {
	return c.JSON(http.StatusOK, h.Store)
}
