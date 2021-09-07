package router

import (
	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/healthcheck"
	"github.com/eiladin/go-simple-startpage/pkg/network"
	"github.com/eiladin/go-simple-startpage/pkg/providers"
	"github.com/eiladin/go-simple-startpage/pkg/status"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, prv *providers.Handlers) {
	healthzService := &healthcheck.Service{Handler: prv.Healthcheck}
	e.GET("/api/healthz", healthzService.Get())

	netHandler := &network.Service{Handler: prv.Network}
	e.GET("/api/network", netHandler.Get)
	e.POST("/api/network", netHandler.Create)

	statusHandler := &status.Service{Handler: prv.Status}
	e.GET("/api/status/:name", statusHandler.Get)

	cfgHandler := &config.Service{Handler: prv.Config}
	e.GET("/api/appconfig", cfgHandler.Get)
}
