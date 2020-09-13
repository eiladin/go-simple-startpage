package router

import (
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite
}

func (suite HandlerSuite) TestNewHandler() {
	app := echo.New()
	AddRoutes(
		app,
		&mockStore{},
		&model.Config{},
	)
	e := []string{}
	for _, r := range app.Routes() {
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
	CreateNetworkFunc func(*model.Network) error
	GetNetworkFunc    func(*model.Network) error
	GetSiteFunc       func(*model.Site) error
}

func (m *mockStore) Ping() error                            { return m.PingFunc() }
func (m *mockStore) CreateNetwork(net *model.Network) error { return m.CreateNetworkFunc(net) }
func (m *mockStore) GetNetwork(net *model.Network) error    { return m.GetNetworkFunc(net) }
func (m *mockStore) GetSite(site *model.Site) error         { return m.GetSiteFunc(site) }
