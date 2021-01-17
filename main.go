package main

import (
	"fmt"
	"log"

	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/eiladin/go-simple-startpage/internal/server"
	"github.com/eiladin/go-simple-startpage/internal/yamlstore"
	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/store"
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
	var store store.Store
	var err error
	if len(c.Filepath) > 0 {
		store, err = yamlstore.New(c.Filepath)
		if err != nil {
			log.Fatal(err)
		}
	} else if len(c.Database.Driver) > 0 {
		store, err = database.New(&c.Database)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Filepath or Database must be set in config.yaml")
	}
	if err = store.Ping(); err != nil {
		log.Fatal(err)
	}

	s := server.New(c, store)
	s.Logger.Fatal(s.Start(fmt.Sprintf(":%d", c.ListenPort)))
}
