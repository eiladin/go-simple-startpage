package api

import (
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/stretchr/testify/suite"
)

type mockStore struct {
	NewFunc           func(*models.Config) (store.Store, error)
	PingFunc          func() error
	CreateNetworkFunc func(*models.Network) error
	GetNetworkFunc    func(*models.Network) error
	GetSiteFunc       func(*models.Site) error
}

func (m *mockStore) New(c *models.Config) (store.Store, error) {
	return m.NewFunc(c)
}

func (m *mockStore) Ping() error {
	return m.PingFunc()
}

func (m *mockStore) CreateNetwork(net *models.Network) error {
	return m.CreateNetworkFunc(net)
}

func (m *mockStore) GetNetwork(net *models.Network) error {
	return m.GetNetworkFunc(net)
}

func (m *mockStore) GetSite(site *models.Site) error {
	return m.GetSiteFunc(site)
}

type HandlerSuite struct {
	suite.Suite
}

func (suite HandlerSuite) TestNewHandler() {
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	h := NewHandler(
		app,
		&mockStore{},
		&models.Config{
			Version: "test-handler-version",
		},
	)
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	suite.Equal("test-handler-version", h.Config.Version)
	suite.Contains(e, "GET /api/network")
	suite.Contains(e, "POST /api/network")
	suite.Contains(e, "GET /api/healthz")
	suite.Contains(e, "GET /api/appconfig")
	suite.Contains(e, "GET /api/status/:id")
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}
