package handler

import "github.com/eiladin/go-simple-startpage/pkg/model"

type mockStore struct {
	CreateNetworkFunc func(*model.Network)
	GetNetworkFunc    func(*model.Network)
	GetSiteFunc       func(*model.Site)
}

func (m *mockStore) CreateNetwork(net *model.Network) {
	m.CreateNetworkFunc(net)
}

func (m *mockStore) GetNetwork(net *model.Network) {
	m.GetNetworkFunc(net)
}

func (m *mockStore) GetSite(site *model.Site) {
	m.GetSiteFunc(site)
}
