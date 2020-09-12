package router

import (
	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/eiladin/go-simple-startpage/pkg/service/config"
	"github.com/eiladin/go-simple-startpage/pkg/service/healthcheck"
	"github.com/eiladin/go-simple-startpage/pkg/service/network"
	"github.com/eiladin/go-simple-startpage/pkg/service/status"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"github.com/pangpanglabs/echoswagger/v2"
)

func AddRoutes(api echoswagger.ApiRoot, store store.Store, cfg *models.Config) {
	config.NewConfigService(cfg).Register(api)
	healthcheck.NewHealthcheckService(cfg, store).Register(api)
	network.NewNetworkService(cfg, store).Register(api)
	status.NewStatusService(cfg, store).Register(api)
}
