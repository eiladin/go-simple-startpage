package router

import (
	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/eiladin/go-simple-startpage/pkg/service/config"
	"github.com/eiladin/go-simple-startpage/pkg/service/healthcheck"
	"github.com/eiladin/go-simple-startpage/pkg/service/network"
	"github.com/eiladin/go-simple-startpage/pkg/service/status"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"github.com/labstack/echo/v4"
)

func AddRoutes(api *echo.Echo, store store.Store, cfg *model.Config) {
	config.NewConfigService(cfg).Register(api)
	healthcheck.NewHealthcheckService(cfg, store).Register(api)
	network.NewNetworkService(cfg, store).Register(api)
	status.NewStatusService(cfg, store).Register(api)
}
