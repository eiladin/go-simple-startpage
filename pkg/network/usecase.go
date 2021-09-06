package network

import (
	"errors"
	"sort"
)

var (
	ErrNotFound = errors.New("Network not found")
)

type INetwork interface {
	Get() (*Network, error)
	Create(*Network) error
}

type repository interface {
	CreateNetwork(*Network) error
	GetNetwork(*Network) error
}

// Compile-time proof of interface implementation.
var _ INetwork = (*service)(nil)

type service struct {
	repo repository
}

func New(repo repository) INetwork {
	return &service{
		repo: repo,
	}
}

func (c *service) Create(net *Network) error {
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

func (c *service) Get() (*Network, error) {
	var net Network

	if err := c.repo.GetNetwork(&net); err != nil {
		return nil, ErrNotFound
	}
	sortSitesByName(net.Sites)
	return &net, nil
}
