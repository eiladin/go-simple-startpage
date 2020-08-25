package store

import (
	"errors"

	"github.com/eiladin/go-simple-startpage/pkg/model"
)

// Store interface
type Store interface {
	New() (Store, error)
	CreateNetwork(net *model.Network) error
	GetNetwork(net *model.Network) error
	GetSite(site *model.Site) error
}

// ErrNotFound not found error
var ErrNotFound = errors.New("record not found")
