package router

import (
	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/healthcheck"
	"github.com/eiladin/go-simple-startpage/pkg/network"
	"github.com/eiladin/go-simple-startpage/pkg/providers"
	"github.com/eiladin/go-simple-startpage/pkg/status"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, prv *providers.Provider) {
	healthzHandler := &healthcheck.Handler{UseCase: prv.Healthcheck}
	e.GET("/api/healthz", healthzHandler.Get())

	netHandler := &network.Handler{UseCase: prv.Network}
	e.GET("/api/network", netHandler.Get)
	e.POST("/api/network", netHandler.Create)

	statusHandler := &status.Handler{UseCase: prv.Status}
	e.GET("/api/status/:name", statusHandler.Get)

	cfgHandler := &config.Handler{UseCase: prv.Config}
	e.GET("/api/appconfig", cfgHandler.Get)
}
