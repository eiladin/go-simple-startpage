package handlers

import (
	"errors"
	"net/http"

	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/eiladin/go-simple-startpage/pkg/usecases/network"
	"github.com/labstack/echo/v4"
)

type NetworkHandler struct {
	NetworkUseCase network.INetwork
}

// Create godoc
// @Summary Add Network
// @Description add or update network
// @Tags Network
// @Accept  json
// @Produce  json
// @Param network body models.Network true "Add Network"
// @Success 201 {object} models.NetworkID
// @Failure 400 {object} httperror.HTTPError
// @Failure 500 {object} httperror.HTTPError
// @Router /api/network [post]
func (c *NetworkHandler) Create(ctx echo.Context) error {
	net := new(models.Network)

	if err := ctx.Bind(net); err != nil || (net.Network == "" && net.Links == nil && net.Sites == nil) {
		if err == nil {
			err = errors.New("empty request recieved")
		}
		return echo.ErrBadRequest.SetInternal(err)
	}

	if err := c.NetworkUseCase.Create(net); err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return ctx.JSON(http.StatusCreated, models.NetworkID{ID: net.ID})
}

// Get godoc
// @Summary Get Network
// @Description get network
// @Tags Network
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Network
// @Failure 404 {object} httperror.HTTPError
// @Failure 503 {object} httperror.HTTPError
// @Router /api/network [get]
func (c *NetworkHandler) Get(ctx echo.Context) error {
	net, err := c.NetworkUseCase.Get()

	if err != nil {
		if errors.Is(err, network.ErrNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return ctx.JSON(http.StatusOK, net)
}
