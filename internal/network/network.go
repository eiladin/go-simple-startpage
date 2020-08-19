package network

import (
	"net/http"
	"sort"

	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
)

// Store interface
type Store interface {
	CreateNetwork(net *model.Network)
	GetNetwork(net *model.Network)
}

// Handler handles Network commands
type Handler struct {
	Store Store
}

// Get handles /api/network
func (h Handler) Get(c echo.Context) error {
	var net model.Network
	h.Store.GetNetwork(&net)
	sort.Slice(net.Sites, func(p, q int) bool {
		return net.Sites[p].FriendlyName < net.Sites[q].FriendlyName
	})
	return c.JSON(http.StatusOK, net)
}

// Create handles /api/network
func (h Handler) Create(c echo.Context) error {
	net := new(model.Network)
	err := c.Bind(net)
	if err != nil || (net.Network == "" && net.ID == 0 && net.Links == nil && net.Sites == nil) {
		return echo.ErrBadRequest
	}

	h.Store.CreateNetwork(net)

	return c.JSON(http.StatusCreated, model.NetworkID{ID: net.ID})
}
