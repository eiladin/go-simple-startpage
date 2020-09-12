package router

import (
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
}

func (suite HandlerSuite) TestNewHandler() {
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	AddRoutes(
		app,
		&mockStore{},
		&models.Config{},
	)
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	suite.Contains(e, "GET /api/network")
	suite.Contains(e, "POST /api/network")
	suite.Contains(e, "GET /api/healthz")
	suite.Contains(e, "GET /api/appconfig")
	suite.Contains(e, "GET /api/status/:id")
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

type mockStore struct {
	mock.Mock
	PingFunc          func() error
	CreateNetworkFunc func(*models.Network) error
	GetNetworkFunc    func(*models.Network) error
	GetSiteFunc       func(*models.Site) error
}

func (m *mockStore) Ping() error                             { return m.PingFunc() }
func (m *mockStore) CreateNetwork(net *models.Network) error { return m.CreateNetworkFunc(net) }
func (m *mockStore) GetNetwork(net *models.Network) error    { return m.GetNetworkFunc(net) }
func (m *mockStore) GetSite(site *models.Site) error         { return m.GetSiteFunc(site) }
