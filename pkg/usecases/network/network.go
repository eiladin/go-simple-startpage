package network

import (
	"errors"
	"sort"

	"github.com/eiladin/go-simple-startpage/pkg/models"
)

var (
	ErrNotFound = errors.New("Network not found")
)

type INetwork interface {
	Get() (*models.Network, error)
	Create(*models.Network) error
}

type repository interface {
	CreateNetwork(*models.Network) error
	GetNetwork(*models.Network) error
}

type service struct {
	repo repository
}

func New(repo repository) INetwork {
	return &service{
		repo: repo,
	}
}

func (c *service) Create(net *models.Network) error {
	if err := c.repo.CreateNetwork(net); err != nil {
		return err
	}
	return nil
}

func sortSitesByName(sites []models.Site) {
	sort.Slice(sites, func(p, q int) bool {
		return sites[p].FriendlyName < sites[q].FriendlyName
	})
}

func (c *service) Get() (*models.Network, error) {
	var net models.Network

	if err := c.repo.GetNetwork(&net); err != nil {
		return nil, ErrNotFound
	}
	sortSitesByName(net.Sites)
	return &net, nil
}
