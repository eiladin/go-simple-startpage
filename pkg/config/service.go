package config

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Service struct {
	Handler IHandler
}

// Get godoc
// @Summary Get AppConfig
// @Description get application configuration
// @Tags AppConfig
// @Accept  json
// @Produce  json
// @Success 200 {object} config.Config
// @Router /api/appconfig [get]
func (c Service) Get(ctx echo.Context) error {
	cfg, err := c.Handler.Get()
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}
	return ctx.JSON(http.StatusOK, cfg)
}
