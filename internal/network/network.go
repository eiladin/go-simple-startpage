package network

import (
	"net/http"

	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
	"github.com/labstack/echo/v4"
)

// Handler handles Network commands
type Handler struct {
	NetworkService interfaces.NetworkService
}

// Get handles /api/network
func (h Handler) Get(c echo.Context) error {
	var net interfaces.Network
	h.NetworkService.FindNetwork(&net)
	return c.JSON(http.StatusOK, net)
}

// Create handles /api/network
func (h Handler) Create(c echo.Context) error {
	net := new(interfaces.Network)
	err := c.Bind(net)
	if err != nil || (net.Network == "" && net.ID == 0 && net.Links == nil && net.Sites == nil) {
		return echo.ErrBadRequest
	}

	h.NetworkService.CreateNetwork(net)

	return c.JSON(http.StatusCreated, interfaces.NetworkID{ID: net.ID})
}
