package store

import (
	"errors"

	"github.com/eiladin/go-simple-startpage/pkg/model"
)

var ErrNotFound = errors.New("record not found")

type Store interface {
	Ping() error
	CreateNetwork(net *model.Network) error
	GetNetwork(net *model.Network) error
	GetSite(site *model.Site) error
}
