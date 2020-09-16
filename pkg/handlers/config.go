package handlers

import (
	"net/http"

	"github.com/eiladin/go-simple-startpage/pkg/usecases/config"
	"github.com/labstack/echo/v4"
)

type ConfigHandler struct {
	ConfigUseCase config.IConfig
}

// Get godoc
// @Summary Get AppConfig
// @Description get application configuration
// @Tags AppConfig
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Config
// @Router /api/appconfig [get]
func (c ConfigHandler) Get(ctx echo.Context) error {
	cfg, err := c.ConfigUseCase.Get()
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}
	return ctx.JSON(http.StatusOK, cfg)
}
