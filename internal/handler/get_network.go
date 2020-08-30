package handler

import (
	"errors"
	"net/http"
	"sort"

	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

func (h handler) getNetwork(c echo.Context) error {
	var net models.Network
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

func (h handler) AddGetNetworkRoute(app echoswagger.ApiRoot) echoswagger.ApiRoot {
	app.GET("/api/network", h.getNetwork).
		AddResponse(http.StatusOK, "success", models.Network{}, nil).
		AddResponse(http.StatusNotFound, "not found", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)

	return app
}
