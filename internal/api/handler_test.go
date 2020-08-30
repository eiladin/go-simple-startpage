package api

import (
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/stretchr/testify/assert"
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

func newMockStore() mockStore {
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
	return s
}

func newMockHandler() handler {
	s := newMockStore()
	return handler{Store: &s}
}

func TestGetHandler(t *testing.T) {
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	h := NewHandler(
		app,
		newMockStore(),
		&models.Config{
			Version: "test-handler-version",
		},
	)
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	assert.Equal(t, "test-handler-version", h.Config.Version)
	assert.Contains(t, e, "GET /api/network")
	assert.Contains(t, e, "POST /api/network")
	assert.Contains(t, e, "GET /api/healthz")
	assert.Contains(t, e, "GET /api/appconfig")
	assert.Contains(t, e, "GET /api/status/:id")
}
