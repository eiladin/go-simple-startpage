package main

import (
	"fmt"
	"log"

	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/eiladin/go-simple-startpage/internal/server"
	"github.com/eiladin/go-simple-startpage/pkg/config"
)

var version = "dev"

// @title Go Simple Startpage API
// @description This is the API for the Go Simple Startpage App

// @contact.name Sami Khan
// @contact.url https://github.com/eiladin/go-simple-startpage

// @license.name MIT
// @license.url https://github.com/eiladin/go-simple-startpage/blob/master/LICENSE
//go:generate swag init -o ./internal/server/docs
func main() {
	c := config.Load(version, "")
	store, err := database.New(&c.Database)
	if err != nil {
		log.Fatal(err)
	}
	if err = store.Ping(); err != nil {
		log.Fatal(err)
	}

	s := server.New(c, store)
	s.Logger.Fatal(s.Start(fmt.Sprintf(":%d", c.ListenPort)))
}
