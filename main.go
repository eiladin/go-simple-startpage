package main

import (
	"github.com/eiladin/go-simple-startpage/config"
	"github.com/eiladin/go-simple-startpage/db"
	"github.com/eiladin/go-simple-startpage/handlers"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	config := config.InitConfig()

	db.InitDB()
	db.MigrateDB()

	r := gin.Default()
	r.Use(CORSMiddleware())
	r.Use(static.Serve("/", static.LocalFile("./ui/dist/ui", true)))
	r.NoRoute(func(c *gin.Context) {
		c.File("./ui/dist/ui/index.html")
	})

	r.GET("/api/appconfig", handlers.GetConfigHandler)
	r.GET("/api/network", handlers.GetNetworkHandler)
	r.POST("/api/network", handlers.AddNetworkHandler)
	r.GET("/api/status", handlers.GetStatusHandler)
	r.POST("/api/status", handlers.UpdateStatusHandler)

	err := r.Run(":" + config.Server.Port)
	if err != nil {
		panic(err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, GET, OPTIONS, POST, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
