package status

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Service struct {
	Handler IHandler
}

// Get godoc
// @Summary Get Status
// @Description get status given a site id
// @Tags Status
// @Accept  json
// @Produce  json
// @Param  id path int true "Site ID"
// @Success 200 {object} Status
// @Failure 400 {object} httperror.HTTPError
// @Failure 404 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /api/status/{name} [get]
func (c *Service) Get(ctx echo.Context) error {
	name := ctx.Param("name")
	if name == "" {
		return echo.ErrBadRequest
	}

	s, err := c.Handler.Get(name)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return ctx.JSON(http.StatusOK, s)
}
