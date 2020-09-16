package store

import "github.com/eiladin/go-simple-startpage/pkg/models"

type Store interface {
	Ping() error
	CreateNetwork(net *models.Network) error
	GetNetwork(net *models.Network) error
	GetSite(site *models.Site) error
}
