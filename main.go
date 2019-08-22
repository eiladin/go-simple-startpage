package main

import (
	"github.com/eiladin/go-simple-startpage/handlers"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.Use(static.Serve("/", static.LocalFile("./ui/dist/ui", true)))
	r.NoRoute(func(c *gin.Context) {
		c.File("./ui/dist/ui/index.html")
	})

	r.GET("/api/config", handlers.GetConfigHandler)
	r.POST("/api/config", handlers.AddConfigHandler)

	err := r.Run(":3000")
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
