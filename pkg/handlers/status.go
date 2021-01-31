package handlers

import (
	"errors"
	"net/http"

	"github.com/eiladin/go-simple-startpage/pkg/usecases/status"
	"github.com/labstack/echo/v4"
)

type StatusHandler struct {
	StatusUseCase status.IStatus
}

// Get godoc
// @Summary Get Status
// @Description get status given a site id
// @Tags Status
// @Accept  json
// @Produce  json
// @Param  id path int true "Site ID"
// @Success 200 {object} models.Status
// @Failure 400 {object} httperror.HTTPError
// @Failure 404 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /api/status/{name} [get]
func (c *StatusHandler) Get(ctx echo.Context) error {
	// httpClient.Timeout = time.Millisecond * time.Duration(c.config.Timeout)
	name := ctx.Param("name")
	if name == "" {
		return echo.ErrBadRequest
	}

	s, err := c.StatusUseCase.Get(name)
	if err != nil {
		if errors.Is(err, status.ErrNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return ctx.JSON(http.StatusOK, s)
}
