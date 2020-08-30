package store

import (
	"errors"

	"github.com/eiladin/go-simple-startpage/internal/models"
)

var ErrNotFound = errors.New("record not found")

type Store interface {
	New(config *models.Config) (Store, error)
	CreateNetwork(net *models.Network) error
	GetNetwork(net *models.Network) error
	GetSite(site *models.Site) error
}
