package main

import (
	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/eiladin/go-simple-startpage/internal/network"
	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func setupMiddleware(app *echo.Echo) {
	app.Use(middleware.CORS())
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	app.Use(middleware.Recover())
	app.Use(middleware.Gzip())
}

func setupRoutes(app *echo.Echo, store *database.DB) {
	handler := network.Handler{NetworkService: store}
	app.GET("/api/appconfig", config.GetAppConfig)
	app.GET("/api/network", handler.GetNetwork)
	app.POST("/api/network", handler.NewNetwork)
	app.GET("/api/status/:id", handler.GetStatus)
	app.Static("/", "./ui/dist/ui")
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

	app := echo.New()
	setupMiddleware(app)
	setupRoutes(app, &store)

	app.Logger.Fatal(app.Start(":" + c.Server.Port))
}
