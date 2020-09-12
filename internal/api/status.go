package api

import (
	"crypto/tls"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

type StatusService struct {
	config *models.Config
	store  store.Store
}

func NewStatusService(cfg *models.Config, store store.Store) StatusService {
	return StatusService{
		config: cfg,
		store:  store,
	}
}

var httpClient = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func (s StatusService) Get(ctx echo.Context) error {
	httpClient.Timeout = time.Millisecond * time.Duration(s.config.Timeout)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id < 1 {
		if err == nil {
			err = errors.New("invalid id received: " + ctx.Param("id"))
		}
		return echo.ErrBadRequest.SetInternal(err)
	}

	site := models.Site{ID: uint(id)}
	err = s.store.GetSite(&site)
	if err != nil {
		return echo.ErrNotFound
	}

	res := models.NewSiteStatus(httpClient, &site)
	return ctx.JSON(http.StatusOK, res)
}

func (s StatusService) Register(api echoswagger.ApiRoot) {
	api.GET("/api/status/:id", s.Get).
		AddParamPath(0, "id", "SiteID to get status for").
		AddResponse(http.StatusOK, "success", models.SiteStatus{}, nil).
		AddResponse(http.StatusBadRequest, "bad request", nil, nil).
		AddResponse(http.StatusNotFound, "not found", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)
}
