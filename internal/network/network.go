package network

import (
	"fmt"
	"net/http"

	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
	"github.com/labstack/echo"
)

// Handler handles Network commands
type Handler struct {
	NetworkService interfaces.NetworkService
}

// GetNetwork handles /api/network
func (h Handler) GetNetwork(c echo.Context) error {
	var net interfaces.Network
	h.NetworkService.FindNetwork(&net)
	return c.JSON(http.StatusOK, net)
}

// NewNetwork handles /api/network
func (h Handler) NewNetwork(c echo.Context) error {
	net := new(interfaces.Network)
	err := c.Bind(net)
	if err != nil {
		return echo.ErrBadRequest
	}

	h.NetworkService.CreateNetwork(net)

	return c.String(http.StatusCreated, fmt.Sprintf(`{"id":"%d"}`, net.ID))
}
