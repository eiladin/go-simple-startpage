package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/etherlabsio/healthcheck"
	"github.com/labstack/echo/v4"
)

func (h handler) checkDB(ctx context.Context) error {
	err := h.Store.Ping()
	if err != nil {
		return fmt.Errorf("unable to connect to database %w", err)
	}
	return nil
}

func (h handler) getHeathcheck() echo.HandlerFunc {
	return echo.WrapHandler(healthcheck.Handler(
		healthcheck.WithTimeout(5*time.Second),

		healthcheck.WithChecker("database", healthcheck.CheckerFunc(h.checkDB)),
	))
}

func (h handler) addHeathcheckRoutes() {
	h.GET("/api/healthz", h.getHeathcheck()).
		AddResponse(http.StatusOK, "success", models.Healthcheck{}, nil)
}
