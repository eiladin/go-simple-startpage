package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/etherlabsio/healthcheck"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

type HealthcheckService struct {
	config *models.Config
	store  store.Store
}

func NewHealthcheckService(cfg *models.Config, store store.Store) HealthcheckService {
	return HealthcheckService{
		config: cfg,
		store:  store,
	}
}

func (s HealthcheckService) checkDB(ctx context.Context) error {
	err := s.store.Ping()
	if err != nil {
		return fmt.Errorf("unable to connect to database %w", err)
	}
	return nil
}

func (s HealthcheckService) Get() echo.HandlerFunc {
	return echo.WrapHandler(healthcheck.Handler(
		healthcheck.WithTimeout(5*time.Second),

		healthcheck.WithChecker("database", healthcheck.CheckerFunc(s.checkDB)),
	))
}

func (s HealthcheckService) Register(api echoswagger.ApiRoot) {
	api.GET("/api/healthz", s.Get()).
		AddResponse(http.StatusOK, "success", models.Healthcheck{}, nil)
}
