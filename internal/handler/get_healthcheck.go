package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/etherlabsio/healthcheck"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

func (h handler) checkDB(ctx context.Context) error {
	_, err := (&database.DB{}).New(h.Config)
	if err != nil {
		return fmt.Errorf("unable to connect to database %w", err)
	}
	return nil
}

func (h handler) GetHeathcheck() http.Handler {
	return healthcheck.Handler(
		healthcheck.WithTimeout(5*time.Second),

		healthcheck.WithChecker("database", healthcheck.CheckerFunc(h.checkDB)),
	)
}

func (h handler) AddGetHealthzRoute(app echoswagger.ApiRoot) echoswagger.ApiRoot {
	app.GET("/api/healthz", echo.WrapHandler(h.GetHeathcheck())).
		AddResponse(http.StatusOK, "success", models.Healthcheck{}, nil)

	return app
}
