package handlers

import (
	"github.com/eiladin/go-simple-startpage/pkg/providers"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, prv *providers.Provider) {
	healthzHandler := &HealthcheckHandler{prv.Healthcheck}
	e.GET("/api/healthz", healthzHandler.Get())

	netHandler := &NetworkHandler{prv.Network}
	e.GET("/api/network", netHandler.Get)
	e.POST("/api/network", netHandler.Create)

	statusHandler := &StatusHandler{prv.Status}
	e.GET("/api/status/:id", statusHandler.Get)

	cfgHandler := &ConfigHandler{prv.Config}
	e.GET("/api/appconfig", cfgHandler.Get)
}
