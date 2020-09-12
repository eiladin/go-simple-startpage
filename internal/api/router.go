package api

import (
	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/pangpanglabs/echoswagger/v2"
)

func AddRoutes(api echoswagger.ApiRoot, store store.Store, config *models.Config) {
	NewConfigService(config).Register(api)
	NewHealthcheckService(config, store).Register(api)
	NewNetworkService(config, store).Register(api)
	NewStatusService(config, store).Register(api)
}
