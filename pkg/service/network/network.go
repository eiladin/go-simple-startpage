package network

import (
	"errors"
	"net/http"
	"sort"

	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

type NetworkService struct {
	config *models.Config
	store  store.Store
}

func NewNetworkService(cfg *models.Config, store store.Store) NetworkService {
	return NetworkService{
		config: cfg,
		store:  store,
	}
}

func (s NetworkService) Create(ctx echo.Context) error {
	net := new(models.Network)

	if err := ctx.Bind(net); err != nil || (net.Network == "" && net.ID == 0 && net.Links == nil && net.Sites == nil) {
		if err == nil {
			err = errors.New("empty request recieved")
		}
		return echo.ErrBadRequest.SetInternal(err)
	}

	if err := s.store.CreateNetwork(net); err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return ctx.JSON(http.StatusCreated, models.NetworkID{ID: net.ID})
}

func sortSitesByName(sites []models.Site) {
	sort.Slice(sites, func(p, q int) bool {
		return sites[p].FriendlyName < sites[q].FriendlyName
	})
}

func (s NetworkService) Get(ctx echo.Context) error {
	var net models.Network

	if err := s.store.GetNetwork(&net); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}
	sortSitesByName(net.Sites)

	return ctx.JSON(http.StatusOK, net)
}

func (s NetworkService) Register(api echoswagger.ApiRoot) {
	api.POST("/api/network", s.Create).
		AddParamBody(models.Network{}, "body", "Network to add", true).
		AddResponse(http.StatusCreated, "success", models.NetworkID{}, nil).
		AddResponse(http.StatusBadRequest, "bad request", nil, nil).
		AddResponse(http.StatusNotFound, "not found", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)

	api.GET("/api/network", s.Get).
		AddResponse(http.StatusOK, "success", models.Network{}, nil).
		AddResponse(http.StatusNotFound, "not found", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)
}
