package providers

import (
	"github.com/eiladin/go-simple-startpage/pkg/config"
	cfg "github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/healthcheck"
	"github.com/eiladin/go-simple-startpage/pkg/network"
	"github.com/eiladin/go-simple-startpage/pkg/status"
	"github.com/eiladin/go-simple-startpage/pkg/store"
)

type Handlers struct {
	Network     network.IHandler
	Healthcheck healthcheck.IHandler
	Status      status.IHandler
	Config      config.IHandler
}

func InitProvider(cfg *cfg.Config, store store.Store) *Handlers {
	net := network.New(store)
	hc := healthcheck.New(store)
	status := status.New(store, cfg)
	ch := config.New(cfg)

	return &Handlers{
		Network:     net,
		Healthcheck: hc,
		Status:      status,
		Config:      ch,
	}
}
