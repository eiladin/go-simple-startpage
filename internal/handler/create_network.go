package handler

import (
	"errors"
	"net/http"

	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

func (h handler) CreateNetwork(c echo.Context) error {
	net := new(models.Network)
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

	return c.JSON(http.StatusCreated, models.NetworkID{ID: net.ID})
}

func (h handler) AddPostNetworkRoute(app echoswagger.ApiRoot) echoswagger.ApiRoot {
	app.POST("/api/network", h.CreateNetwork).
		AddParamBody(models.Network{}, "body", "Network to add", true).
		AddResponse(http.StatusCreated, "success", models.NetworkID{}, nil).
		AddResponse(http.StatusBadRequest, "bad request", nil, nil).
		AddResponse(http.StatusNotFound, "not found", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)

	return app
}
