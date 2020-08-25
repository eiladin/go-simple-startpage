package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/eiladin/go-simple-startpage/internal/handler"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pangpanglabs/echoswagger/v2"
)

func swaggerRefSkipper(ctx echo.Context) bool {
	return strings.Contains(ctx.Request().Header.Get("Referer"), "swagger")
}

func apiSkipper(ctx echo.Context) bool {
	return strings.Contains(ctx.Request().Header.Get("Referer"), "swagger")
}

func setupMiddleware(app *echo.Echo, c config.Config) {
	if c.IsProduction() {
		app.Use(middleware.Logger())
	} else {
		app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}, error=${error}\n",
		}))
	}

	app.Use(middleware.CORS())
	app.Use(middleware.RequestID())
	app.Use(middleware.Secure())
	app.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{Skipper: swaggerRefSkipper}))
	app.Use(middleware.Recover())
	app.Use(middleware.Gzip())
	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: apiSkipper,
		Index:   "index.html",
		Root:    "ui/dist/ui",
		Browse:  false,
		HTML5:   true,
	}))
}

func setupRoutes(app echoswagger.ApiRoot, store store.Store) {
	app.GET("/api/appconfig", handler.Config{Store: config.GetConfig()}.Get).
		AddResponse(http.StatusOK, "success", config.Config{}, nil)

	app.GET("/api/network", handler.Network{Store: store}.Get).
		AddResponse(http.StatusOK, "success", model.Network{}, nil).
		AddResponse(http.StatusNotFound, "not found", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)

	app.POST("/api/network", handler.Network{Store: store}.Create).
		AddParamBody(model.Network{}, "body", "Network to add", true).
		AddResponse(http.StatusCreated, "success", model.NetworkID{}, nil).
		AddResponse(http.StatusBadRequest, "bad request", nil, nil).
		AddResponse(http.StatusNotFound, "not found", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)

	app.GET("/api/status/:id", handler.Status{Store: store}.Get).
		AddParamPath(0, "id", "SiteID to get status for").
		AddResponse(http.StatusOK, "success", model.SiteStatus{}, nil).
		AddResponse(http.StatusBadRequest, "bad request", nil, nil).
		AddResponse(http.StatusNotFound, "not found", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)
}

var version = "dev"

func main() {
	c := config.InitConfig(version, "")
	store, err := database.DB{}.New()
	if err != nil {
		log.Fatal(err)
	}

	app := echoswagger.New(echo.New(), "/swagger", &echoswagger.Info{
		Title:       "Go Simple Startpage API",
		Description: "This in the API for the Go Simple Startpage App",
		Version:     "1.0.0",
		Contact: &echoswagger.Contact{
			Name: "Sami Khan",
			URL:  "https://github.com/eiladin/go-simple-startpage",
		},
		License: &echoswagger.License{
			Name: "MIT",
			URL:  "https://github.com/eiladin/go-simple-startpage/blob/master/LICENSE",
		},
	})

	e := app.Echo()

	setupMiddleware(e, c)
	setupRoutes(app, store)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", c.ListenPort)))
}
