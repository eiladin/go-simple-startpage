package handler

import (
	"errors"
	"net/http"
	"sort"

	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

// Network struct
type Network struct {
	Store store.Store
}

// Get /api/network
func (h Network) Get(c echo.Context) error {
	var net model.Network
	err := h.Store.GetNetwork(&net)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}
	sort.Slice(net.Sites, func(p, q int) bool {
		return net.Sites[p].FriendlyName < net.Sites[q].FriendlyName
	})
	return c.JSON(http.StatusOK, net)
}

// Create /api/network
func (h Network) Create(c echo.Context) error {
	net := new(model.Network)
	err := c.Bind(net)
	if err != nil || (net.Network == "" && net.ID == 0 && net.Links == nil && net.Sites == nil) {
		if err == nil {
			err = errors.New("empty request recieved")
		}
		return echo.ErrBadRequest.SetInternal(err)
	}

	err = h.Store.CreateNetwork(net)
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return c.JSON(http.StatusCreated, model.NetworkID{ID: net.ID})
}

// Register handler
func (h Network) Register(app echoswagger.ApiRoot) echoswagger.ApiRoot {
	app.GET("/api/network", h.Get).
		AddResponse(http.StatusOK, "success", model.Network{}, nil).
		AddResponse(http.StatusNotFound, "not found", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)

	app.POST("/api/network", h.Create).
		AddParamBody(model.Network{}, "body", "Network to add", true).
		AddResponse(http.StatusCreated, "success", model.NetworkID{}, nil).
		AddResponse(http.StatusBadRequest, "bad request", nil, nil).
		AddResponse(http.StatusNotFound, "not found", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)

	return app
}
