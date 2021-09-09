package router

import (
	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/healthcheck"
	"github.com/eiladin/go-simple-startpage/pkg/network"
	"github.com/eiladin/go-simple-startpage/pkg/status"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"github.com/labstack/echo/v4"
)

type Handlers struct {
	Network     network.IHandler
	Healthcheck healthcheck.IHandler
	Status      status.IHandler
	Config      config.IHandler
}

func DefaultHandlers(cfg *config.Config, store store.Store) *Handlers {
	net := network.New(store)
	hc := healthcheck.New(store)
	status := status.New(store, cfg)
	ch := config.New(cfg)

	return &Handlers{
		Network:     net,
		Healthcheck: hc,
		Status:      status,
		Config:      ch,
	}
}

func RegisterRoutes(e *echo.Echo, handlers *Handlers) {
	healthzService := &healthcheck.Service{Handler: handlers.Healthcheck}
	e.GET("/api/healthz", healthzService.Get())

	netHandler := &network.Service{Handler: handlers.Network}
	e.GET("/api/network", netHandler.Get)
	e.POST("/api/network", netHandler.Create)

	statusHandler := &status.Service{Handler: handlers.Status}
	e.GET("/api/status/:name", statusHandler.Get)

	cfgHandler := &config.Service{Handler: handlers.Config}
	e.GET("/api/appconfig", cfgHandler.Get)
}
