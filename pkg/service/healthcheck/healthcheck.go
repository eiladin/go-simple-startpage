package healthcheck

import (
	"context"
	"fmt"
	"time"

	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"github.com/etherlabsio/healthcheck"
	"github.com/labstack/echo/v4"
)

type HealthcheckService struct {
	config *model.Config
	store  store.Store
}

func NewHealthcheckService(cfg *model.Config, store store.Store) HealthcheckService {
	return HealthcheckService{config: cfg, store: store}
}

func (s HealthcheckService) checkDB(ctx context.Context) error {
	if err := s.store.Ping(); err != nil {
		return fmt.Errorf("unable to connect to database %w", err)
	}
	return nil
}

// Get godoc
// @Summary Get Health
// @Description run healthcheck
// @Tags HealthCheck
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Healthcheck
// @Failure 503 {object} model.Healthcheck
// @Router /api/healthz [get]
func (s HealthcheckService) Get() echo.HandlerFunc {
	return echo.WrapHandler(healthcheck.Handler(
		healthcheck.WithTimeout(5*time.Second),

		healthcheck.WithChecker("database", healthcheck.CheckerFunc(s.checkDB)),
	))
}

func (s HealthcheckService) Register(api *echo.Echo) {
	api.GET("/api/healthz", s.Get())
}
