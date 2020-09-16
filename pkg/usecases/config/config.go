package config

import (
	"github.com/eiladin/go-simple-startpage/pkg/config"
)

type IConfig interface {
	Get() (*config.Config, error)
}

type service struct {
	config *config.Config
}

func New(cfg *config.Config) IConfig {
	return &service{
		config: cfg,
	}
}

func (c *service) Get() (*config.Config, error) {
	return c.config, nil
}
