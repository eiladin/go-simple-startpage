package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/eiladin/go-simple-startpage/internal/api"
	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/eiladin/go-simple-startpage/internal/models"
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

func setupMiddleware(app *echo.Echo, c *models.Config) {
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
		Root:    "ui/dist",
		Browse:  false,
		HTML5:   true,
	}))
}

var version = "dev"

func main() {
	c := models.NewConfig(version, "")
	store, err := database.DB{}.New(c)
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
	api.NewHandler(app, store, c)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", c.ListenPort)))
}
