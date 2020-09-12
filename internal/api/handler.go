package api

import (
	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/pangpanglabs/echoswagger/v2"
)

type handler struct {
	echoswagger.ApiRoot
	Config *models.Config
	Store  store.Store
}

func NewHandler(app echoswagger.ApiRoot, store store.Store, config *models.Config) handler {
	NewConfigService(config).Register(app)
	h := handler{ApiRoot: app, Config: config, Store: store}
	h.addHeathcheckRoutes()
	h.addNetworkRoutes()
	h.addStatusRoutes()
	return h
}
