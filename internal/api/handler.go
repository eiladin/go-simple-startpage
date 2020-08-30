package api

import (
	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/pangpanglabs/echoswagger/v2"
)

type handler struct {
	Config *models.Config
	Store  store.Store
}

func NewHandler(app echoswagger.ApiRoot, store store.Store, config *models.Config) *handler {
	h := &handler{Config: config, Store: store}
	h.addGetHealthcheckRoute(app)
	h.addGetConfigRoute(app)
	h.addCreateNetworkRoute(app)
	h.addGetNetworkRoute(app)
	h.addGetStatusRoute(app)
	return h
}
