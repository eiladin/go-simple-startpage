package main

import (
	"fmt"
	"log"

	"github.com/eiladin/go-simple-startpage/docs"
	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/eiladin/go-simple-startpage/internal/middleware"
	"github.com/eiladin/go-simple-startpage/internal/router"
	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

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

	e.Use(middleware.GetMiddleware(c)...)
	router.AddRoutes(e, store, c)
	if c.IsProduction() {
		e.GET("/swagger/doc.json", echoSwagger.WrapHandler)
	} else {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", c.ListenPort)))
}
