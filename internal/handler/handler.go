package handler

import (
	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/eiladin/go-simple-startpage/internal/store"
)

type handler struct {
	Config *models.Config
	Store  store.Store
}

func NewHandler(store store.Store, config *models.Config) *handler {
	return &handler{Config: config, Store: store}
}
