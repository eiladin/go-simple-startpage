package network

import (
	"errors"
	"sort"
)

var (
	ErrNotFound = errors.New("Network not found")
)

type IHandler interface {
	Get() (*Network, error)
	Create(*Network) error
}

type repository interface {
	CreateNetwork(*Network) error
	GetNetwork(*Network) error
}

// Compile-time proof of interface implementation.
var _ IHandler = (*handler)(nil)

type handler struct {
	repo repository
}

func NewHandler(repo repository) IHandler {
	return &handler{
		repo: repo,
	}
}

func (c *handler) Create(net *Network) error {
	if err := c.repo.CreateNetwork(net); err != nil {
		return err
	}
	return nil
}

func sortSitesByName(sites []Site) {
	sort.Slice(sites, func(p, q int) bool {
		return sites[p].Name < sites[q].Name
	})
}

func (c *handler) Get() (*Network, error) {
	var net Network

	if err := c.repo.GetNetwork(&net); err != nil {
		return nil, ErrNotFound
	}
	sortSitesByName(net.Sites)
	return &net, nil
}
