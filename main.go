package main

import (
	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/eiladin/go-simple-startpage/internal/network"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

func setupMiddleware(app *fiber.App) {
	app.Use(cors.New())
	app.Use(middleware.Compress())
	app.Use(middleware.Logger())
}

func setupRoutes(app *fiber.App) {
	app.Get("/api/appconfig", config.GetAppConfig)
	app.Get("/api/network", network.GetNetwork)
	app.Post("/api/network", network.NewNetwork)
	app.Post("/api/status", network.UpdateStatus)
	app.Static("/", "./ui/dist/ui", fiber.Static{
		Compress: true,
		Browse:   true,
		Index:    "index.html",
	})
	app.Static("*", "./ui/dist/ui/index.html")
}

func initDatabase() {
	database.InitDB()
	db := database.DBConn
	db.AutoMigrate(&network.Network{})
	db.AutoMigrate(&network.Site{})
	db.AutoMigrate(&network.Tag{})
	db.AutoMigrate(&network.Link{})
}

var version = " dev"

func main() {
	c := config.InitConfig(version)
	initDatabase()
	// defer database.DBConn

	app := fiber.New()
	setupMiddleware(app)
	setupRoutes(app)

	err := app.Listen(c.Server.Port)
	if err != nil {
		panic(err)
	}
}
