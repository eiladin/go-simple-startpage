package providers

import (
	cfg "github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"github.com/eiladin/go-simple-startpage/pkg/usecases/config"
	"github.com/eiladin/go-simple-startpage/pkg/usecases/healthcheck"
	"github.com/eiladin/go-simple-startpage/pkg/usecases/network"
	"github.com/eiladin/go-simple-startpage/pkg/usecases/status"
)

type Provider struct {
	Network     network.INetwork
	Healthcheck healthcheck.IHealthcheck
	Status      status.IStatus
	Config      config.IConfig
}

func InitProvider(cfg *cfg.Config, store store.Store) *Provider {
	net := network.New(store)
	hc := healthcheck.New(store)
	status := status.New(store, cfg)
	ch := config.New(cfg)

	return &Provider{
		Network:     net,
		Healthcheck: hc,
		Status:      status,
		Config:      ch,
	}
}
