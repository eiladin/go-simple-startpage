package handler

import (
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/eiladin/go-simple-startpage/pkg/models"
)

type mockStore struct {
	NewFunc           func(*models.Config) (store.Store, error)
	CreateNetworkFunc func(*models.Network) error
	GetNetworkFunc    func(*models.Network) error
	GetSiteFunc       func(*models.Site) error
}

func (m mockStore) New(c *models.Config) (store.Store, error) {
	return m.NewFunc(c)
}

func (m mockStore) CreateNetwork(net *models.Network) error {
	return m.CreateNetworkFunc(net)
}

func (m mockStore) GetNetwork(net *models.Network) error {
	return m.GetNetworkFunc(net)
}

func (m mockStore) GetSite(site *models.Site) error {
	return m.GetSiteFunc(site)
}

func getMockHandler() handler {
	s := mockStore{
		CreateNetworkFunc: func(net *models.Network) error {
			net.ID = 12345
			return nil
		},
		GetNetworkFunc: func(net *models.Network) error {
			net.ID = 12345
			net.Network = "test-network"
			net.Sites = []models.Site{
				{ID: 1, FriendlyName: "z"},
				{ID: 2, FriendlyName: "a"},
			}
			return nil
		},
		GetSiteFunc: func(site *models.Site) error { return nil },
	}

	return handler{Store: &s}
}
