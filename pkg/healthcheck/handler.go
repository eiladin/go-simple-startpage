package healthcheck

import (
	"github.com/labstack/echo/v4"
)

type Handler struct {
	UseCase IHealthcheck
}

// Get godoc
// @Summary Get Health
// @Description run healthcheck
// @Tags HealthCheck
// @Accept  json
// @Produce  json
// @Success 200 {object} Healthcheck
// @Failure 503 {object} Healthcheck
// @Router /api/healthz [get]
func (c *Handler) Get() echo.HandlerFunc {
	return echo.WrapHandler(c.UseCase.Check())
}
