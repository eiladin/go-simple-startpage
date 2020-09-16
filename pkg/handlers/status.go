package handlers

import (
	"errors"
	"net/http"
	"strconv"

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
// @Success 200 {object} models.SiteStatus
// @Failure 400 {object} httperror.HTTPError
// @Failure 404 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /api/status/{id} [get]
func (c *StatusHandler) Get(ctx echo.Context) error {
	// httpClient.Timeout = time.Millisecond * time.Duration(c.config.Timeout)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id < 1 {
		if err == nil {
			err = errors.New("invalid id received: " + ctx.Param("id"))
		}
		return echo.ErrBadRequest.SetInternal(err)
	}

	s, err := c.StatusUseCase.Get(uint(id))
	if err != nil {
		if errors.Is(err, status.ErrNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return ctx.JSON(http.StatusOK, s)
}
