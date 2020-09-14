package config

import (
	"net/http"

	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
)

type ConfigService struct {
	config *model.Config
}

func NewConfigService(cfg *model.Config) ConfigService {
	return ConfigService{
		config: cfg,
	}
}

// Get godoc
// @Summary Get AppConfig
// @Description get application configuration
// @Tags AppConfig
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Config
// @Router /api/appconfig [get]
func (s ConfigService) Get(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, s.config)
}

func (s ConfigService) Register(api *echo.Echo) {
	api.GET("/api/appconfig", s.Get)
}
