package store

import (
	"errors"

	"github.com/eiladin/go-simple-startpage/pkg/network"
)

var ErrNotFound = errors.New("record not found")

type Store interface {
	Ping() error
	CreateNetwork(net *network.Network) error
	GetNetwork(net *network.Network) error
	GetSite(site *network.Site) error
}
