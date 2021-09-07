package network

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Service struct {
	Handler IHandler
}

// Create godoc
// @Summary Add Network
// @Description add or update network
// @Tags Network
// @Accept  json
// @Produce  json
// @Param network body Network true "Add Network"
// @Success 201 {object} NetworkID
// @Failure 400 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /api/network [post]
func (c *Service) Create(ctx echo.Context) error {
	net := new(Network)

	if err := ctx.Bind(net); err != nil || (net.Network == "" && net.Links == nil && net.Sites == nil) {
		if err == nil {
			err = errors.New("empty request recieved")
		}
		return echo.ErrBadRequest.SetInternal(err)
	}

	if err := c.Handler.Create(net); err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return ctx.JSON(http.StatusCreated, NetworkID{ID: net.ID})
}

// Get godoc
// @Summary Get Network
// @Description get network
// @Tags Network
// @Accept  json
// @Produce  json
// @Success 200 {object} Network
// @Failure 404 {object} httperror.HTTPError
// @Failure 503 {object} httperror.HTTPError
// @Router /api/network [get]
func (c *Service) Get(ctx echo.Context) error {
	net, err := c.Handler.Get()

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return ctx.JSON(http.StatusOK, net)
}
