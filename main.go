package main

import (
	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/eiladin/go-simple-startpage/internal/network"
	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

func setupMiddleware(app *fiber.App) {
	app.Use(cors.New())
	app.Use(middleware.Compress())
	app.Use(middleware.Logger())
}

func setupRoutes(app *fiber.App, store *database.DB) {
	app.Get("/api/appconfig", config.GetAppConfig)
	app.Get("/api/network", network.Handler{NetworkService: store}.GetNetwork)
	app.Post("/api/network", network.Handler{NetworkService: store}.NewNetwork)
	app.Post("/api/status", network.Handler{NetworkService: store}.UpdateStatus)
	app.Static("/", "./ui/dist/ui", fiber.Static{
		Compress: true,
		Browse:   true,
		Index:    "index.html",
	})
	app.Static("*", "./ui/dist/ui/index.html")
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
	// defer database.DBConn

	app := fiber.New()
	setupMiddleware(app)
	setupRoutes(app, &store)

	err := app.Listen(c.Server.Port)
	if err != nil {
		panic(err)
	}
}
