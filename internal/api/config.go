package api

import (
	"net/http"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

type ConfigService struct {
	config *models.Config
}

func NewConfigService(cfg *models.Config) ConfigService {
	return ConfigService{
		config: cfg,
	}
}

func (s ConfigService) Get(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, s.config)
}

func (s ConfigService) Register(api echoswagger.ApiRoot) {
	api.GET("/api/appconfig", s.Get).
		AddResponse(http.StatusOK, "success", models.Config{}, nil)
}
