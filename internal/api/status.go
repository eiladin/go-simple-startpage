package api

import (
	"crypto/tls"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/labstack/echo/v4"
)

var httpClient = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func (h handler) getStatus(c echo.Context) error {
	httpClient.Timeout = time.Millisecond * time.Duration(h.Config.Timeout)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		if err == nil {
			err = errors.New("invalid id received: " + c.Param("id"))
		}
		return echo.ErrBadRequest.SetInternal(err)
	}

	site := models.Site{ID: uint(id)}
	err = h.Store.GetSite(&site)
	if err != nil {
		return echo.ErrNotFound
	}

	res := models.NewSiteStatus(httpClient, &site)
	return c.JSON(http.StatusOK, res)
}

func (h handler) addStatusRoutes() {
	h.GET("/api/status/:id", h.getStatus).
		AddParamPath(0, "id", "SiteID to get status for").
		AddResponse(http.StatusOK, "success", models.SiteStatus{}, nil).
		AddResponse(http.StatusBadRequest, "bad request", nil, nil).
		AddResponse(http.StatusNotFound, "not found", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)
}
