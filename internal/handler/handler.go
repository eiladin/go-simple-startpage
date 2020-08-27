package handler

import (
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/eiladin/go-simple-startpage/pkg/models"
)

type handler struct {
	Config *models.Config
	Store  store.Store
}

func NewHandler(store store.Store, config *models.Config) *handler {
	return &handler{Config: config, Store: store}
}
