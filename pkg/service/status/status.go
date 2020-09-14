package status

import (
	"crypto/tls"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"github.com/labstack/echo/v4"
)

type StatusService struct {
	config *model.Config
	store  store.Store
}

func NewStatusService(cfg *model.Config, store store.Store) StatusService {
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

// Get godoc
// @Summary Get Status
// @Description get status given a site id
// @Tags Status
// @Accept  json
// @Produce  json
// @Param  id path int true "Site ID"
// @Success 200 {object} model.SiteStatus
// @Failure 400 {object} model.HTTPError
// @Failure 404 {object} model.HTTPError
// @Router /api/status/{id} [get]
func (s StatusService) Get(ctx echo.Context) error {
	httpClient.Timeout = time.Millisecond * time.Duration(s.config.Timeout)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id < 1 {
		if err == nil {
			err = errors.New("invalid id received: " + ctx.Param("id"))
		}
		return echo.ErrBadRequest.SetInternal(err)
	}

	site := model.Site{ID: uint(id)}
	err = s.store.GetSite(&site)
	if err != nil {
		return echo.ErrNotFound
	}

	res := model.NewSiteStatus(httpClient, &site)
	return ctx.JSON(http.StatusOK, res)
}

func (s StatusService) Register(api *echo.Echo) {
	api.GET("/api/status/:id", s.Get)
}
