package handlers

import (
	"github.com/eiladin/go-simple-startpage/pkg/usecases/healthcheck"
	"github.com/labstack/echo/v4"
)

type HealthcheckHandler struct {
	HealthcheckUseCase healthcheck.IHealthcheck
}

// Get godoc
// @Summary Get Health
// @Description run healthcheck
// @Tags HealthCheck
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Healthcheck
// @Failure 503 {object} models.Healthcheck
// @Router /api/healthz [get]
func (c *HealthcheckHandler) Get() echo.HandlerFunc {
	return echo.WrapHandler(c.HealthcheckUseCase.Check())
}
