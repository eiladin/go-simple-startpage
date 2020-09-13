package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiladin/go-simple-startpage/docs"
	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/eiladin/go-simple-startpage/internal/router"
	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func localhostSkipper(ctx echo.Context) bool {
	return strings.Contains(ctx.Request().Host, "localhost")
}

func fromSwaggerSkipper(ctx echo.Context) bool {
	return strings.Contains(ctx.Request().URL.Path, "/swagger")
}

func setupMiddleware(app *echo.Echo, c *model.Config) {
	if c.IsProduction() {
		app.Use(middleware.Logger())
	} else {
		app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Output: os.Stdout,
			Format: "method=${method}, uri=${uri}, status=${status}, error=${error}\n",
		}))
	}

	app.Use(middleware.CORS())
	app.Use(middleware.RequestID())
	app.Use(middleware.Secure())
	app.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		Skipper:      localhostSkipper,
		CookieSecure: true,
	}))
	app.Use(middleware.Recover())
	app.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: fromSwaggerSkipper,
	}))
	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: fromSwaggerSkipper,
		Index:   "index.html",
		Root:    "ui/dist",
		Browse:  false,
		HTML5:   true,
	}))
}

var version = "dev"

// @title Go Simple Startpage API
// @description This is the API for the Go Simple Startpage App

// @contact.name Sami Khan
// @contact.url https://github.com/eiladin/go-simple-startpage

// @license.name MIT
// @license.url https://github.com/eiladin/go-simple-startpage/blob/master/LICENSE
//go:generate swag init
func main() {
	c := model.NewConfig(version, "")
	store, err := database.New(c)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	docs.SwaggerInfo.Version = version

	setupMiddleware(e, c)
	router.AddRoutes(e, store, c)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", c.ListenPort)))
}
