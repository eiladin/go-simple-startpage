package server

import (
	"github.com/eiladin/go-simple-startpage/internal/server/docs"
	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/providers"
	"github.com/eiladin/go-simple-startpage/pkg/router"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	*echo.Echo
}

func New(c *config.Config, s store.Store) Server {
	e := Server{Echo: echo.New()}
	docs.SwaggerInfo.Version = c.Version

	e.Use(getMiddleware(c)...)
	router.RegisterRoutes(e.Echo, providers.InitProvider(c, s))
	if c.IsProduction() {
		e.GET("/swagger/doc.json", echoSwagger.WrapHandler)
	} else {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}
	return e
}
