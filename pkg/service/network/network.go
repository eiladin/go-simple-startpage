package network

import (
	"errors"
	"net/http"
	"sort"

	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"github.com/labstack/echo/v4"
)

type NetworkService struct {
	config *model.Config
	store  store.Store
}

func NewNetworkService(cfg *model.Config, store store.Store) NetworkService {
	return NetworkService{
		config: cfg,
		store:  store,
	}
}

// Create godoc
// @Summary Add Network
// @Description add or update network
// @Tags Network
// @Accept  json
// @Produce  json
// @Param network body model.Network true "Add Network"
// @Success 201 {object} model.NetworkID
// @Failure 400 {object} model.HTTPError
// @Failure 500 {object} model.HTTPError
// @Router /api/network [post]
func (s NetworkService) Create(ctx echo.Context) error {
	net := new(model.Network)

	if err := ctx.Bind(net); err != nil || (net.Network == "" && net.ID == 0 && net.Links == nil && net.Sites == nil) {
		if err == nil {
			err = errors.New("empty request recieved")
		}
		return echo.ErrBadRequest.SetInternal(err)
	}

	if err := s.store.CreateNetwork(net); err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return ctx.JSON(http.StatusCreated, model.NetworkID{ID: net.ID})
}

func sortSitesByName(sites []model.Site) {
	sort.Slice(sites, func(p, q int) bool {
		return sites[p].FriendlyName < sites[q].FriendlyName
	})
}

// Get godoc
// @Summary Get Network
// @Description get network
// @Tags Network
// @Accept  json
// @Produce  json
// @Success 200 {object} model.Network
// @Failure 404 {object} model.HTTPError
// @Failure 503 {object} model.HTTPError
// @Router /api/network [get]
func (s NetworkService) Get(ctx echo.Context) error {
	var net model.Network

	if err := s.store.GetNetwork(&net); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}
	sortSitesByName(net.Sites)

	return ctx.JSON(http.StatusOK, net)
}

func (s NetworkService) Register(api *echo.Echo) {
	api.POST("/api/network", s.Create)

	api.GET("/api/network", s.Get)
}
