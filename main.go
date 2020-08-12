package main

import (
	"net/http"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/eiladin/go-simple-startpage/internal/network"
	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
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
}

func setupRoutes(app echoswagger.ApiRoot, store *database.DB) {
	handler := network.Handler{NetworkService: store}

	app.GET("/api/appconfig", config.GetAppConfig).
		AddResponse(http.StatusOK, "success", config.Configuration{}, nil)

	app.GET("/api/network", handler.GetNetwork).
		AddResponse(http.StatusOK, "success", interfaces.Network{}, nil).
		AddResponse(http.StatusInternalServerError, "error", nil, nil)

	app.POST("/api/network", handler.NewNetwork).
		AddParamBody(interfaces.Network{}, "body", "Network to add to the store", true).
		AddResponse(http.StatusOK, "success", interfaces.NetworkID{}, nil)

	app.GET("/api/status/:id", handler.GetStatus).
		AddParamPath(0, "id", "ID of site to get status for").
		AddResponse(http.StatusOK, "success", interfaces.SiteStatus{}, nil).
		AddResponse(http.StatusBadRequest, "bad request", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)

	app.Echo().Static("/", "./ui/dist/ui")
}

func initDatabase() database.DB {
	conn := database.InitDB()
	conn.AutoMigrate(&interfaces.Network{})
	conn.AutoMigrate(&interfaces.Site{})
	conn.AutoMigrate(&interfaces.Tag{})
	conn.AutoMigrate(&interfaces.Link{})
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

	app.Echo().Logger.Fatal(app.Echo().Start(":" + c.Server.Port))
}
