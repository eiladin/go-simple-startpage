package server

import (
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/server/docs"
	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockStore struct {
	mock.Mock
}

func (m *mockStore) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockStore) CreateNetwork(net *models.Network) error {
	args := m.Called(net)
	return args.Error(0)
}

func (m *mockStore) GetNetwork(net *models.Network) error {
	args := m.Called(net)
	return args.Error(0)
}

func (m *mockStore) GetSite(site *models.Site) error {
	args := m.Called(site)
	return args.Error(0)
}

type ServerSuite struct {
	suite.Suite
}

func (suite ServerSuite) TestNew() {
	cases := []struct {
		cfg      *config.Config
		expected []string
	}{
		{
			cfg: &config.Config{
				Version:     "devtest",
				Environment: "production",
			},
			expected: []string{
				"GET /api/network",
				"POST /api/network",
				"GET /api/healthz",
				"GET /api/appconfig",
				"GET /api/status/:id",
				"GET /swagger/doc.json",
			}},
		{
			cfg: &config.Config{
				Version:     "prodtest",
				Environment: "dev",
			},
			expected: []string{
				"GET /api/network",
				"POST /api/network",
				"GET /api/healthz",
				"GET /api/appconfig",
				"GET /api/status/:id",
				"GET /swagger/*",
			}},
	}
	for _, c := range cases {
		e := New(c.cfg, &mockStore{})
		ep := []string{}
		for _, r := range e.Routes() {
			ep = append(ep, r.Method+" "+r.Path)
		}
		for _, ex := range c.expected {
			suite.Contains(ep, ex)
		}
		suite.Equal(c.cfg.Version, docs.SwaggerInfo.Version)
	}
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
