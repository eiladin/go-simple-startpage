package status

import (
	"errors"

	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/network"
)

var (
	ErrNotFound = errors.New("Status not found")
)

type IHandler interface {
	Get(string) (*Status, error)
}

type repository interface {
	GetSite(*network.Site) error
}

// Compile-time proof of interface implementation.
var _ IHandler = (*handler)(nil)

type handler struct {
	repo   repository
	config *config.Config
}

func New(repo repository, cfg *config.Config) IHandler {
	return &handler{
		repo:   repo,
		config: cfg,
	}
}

func (c *handler) Get(name string) (*Status, error) {
	site := network.Site{Name: name}

	if err := c.repo.GetSite(&site); err != nil {
		return nil, ErrNotFound
	}

	res := NewStatus(c.config.Timeout, &site)
	return &res, nil
}
