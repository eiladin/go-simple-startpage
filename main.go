package main

import (
	"fmt"
	"net/http"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/eiladin/go-simple-startpage/internal/handler"
	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pangpanglabs/echoswagger/v2"
)

func setupMiddleware(app *echo.Echo) {
	app.Use(middleware.CORS())
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	app.Use(middleware.Recover())
	app.Use(middleware.Gzip())
	app.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Index:  "index.html",
		Root:   "ui/dist/ui",
		Browse: false,
		HTML5:  true,
	}))
}

func setupRoutes(app echoswagger.ApiRoot, store *database.DB) {
	app.GET("/api/appconfig", handler.Config{Store: config.GetConfig()}.Get).
		AddResponse(http.StatusOK, "success", config.Config{}, nil)

	app.GET("/api/network", handler.Network{Store: store}.Get).
		AddResponse(http.StatusOK, "success", model.Network{}, nil).
		AddResponse(http.StatusInternalServerError, "error", nil, nil)

	app.POST("/api/network", handler.Network{Store: store}.Create).
		AddParamBody(model.Network{}, "body", "Network to add to the store", true).
		AddResponse(http.StatusOK, "success", model.NetworkID{}, nil)

	app.GET("/api/status/:id", handler.Status{Store: store}.Get).
		AddParamPath(0, "id", "ID of site to get status for").
		AddResponse(http.StatusOK, "success", model.SiteStatus{}, nil).
		AddResponse(http.StatusBadRequest, "bad request", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)
}

func initDatabase() database.DB {
	conn := database.InitDB()
	database.MigrateDB(conn)
	return database.DB{DB: conn}
}

var version = " dev"

func main() {
	c := config.InitConfig(version, "")
	store := initDatabase()

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
			URL:  "http://github.com/eiladin/go-simple-startpage/LICENSE",
		},
	})
	setupMiddleware(app.Echo())
	setupRoutes(app, &store)

	app.Echo().Logger.Fatal(app.Echo().Start(fmt.Sprintf(":%d", c.ListenPort)))
}
