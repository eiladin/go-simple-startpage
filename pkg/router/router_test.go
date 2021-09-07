package router

import (
	"net/http"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/network"
	"github.com/eiladin/go-simple-startpage/pkg/providers"
	"github.com/eiladin/go-simple-startpage/pkg/status"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RouteSuite struct {
	suite.Suite
}

type mockHealthCheckHandler struct {
	mock.Mock
}

func (m *mockHealthCheckHandler) Check() http.Handler {
	m.Called()
	return nil
}

type mockNetworkHandler struct {
	mock.Mock
}

func (m *mockNetworkHandler) Get() (*network.Network, error) {
	args := m.Called()
	data := args.Get(0).(network.Network)
	return &data, args.Error(1)
}

func (m *mockNetworkHandler) Create(net *network.Network) error {
	args := m.Called(net)
	return args.Error(0)
}

type mockConfigHandler struct {
	mock.Mock
}

func (m *mockConfigHandler) Get() (*config.Config, error) {
	args := m.Called()
	data := args.Get(0).(config.Config)
	return &data, args.Error(1)
}

type mockStatusHandler struct {
	mock.Mock
}

func (m *mockStatusHandler) Get(name string) (*status.Status, error) {
	args := m.Called(name)
	data := args.Get(0).(status.Status)
	return &data, args.Error(1)
}

type HandlerSuite struct {
	suite.Suite
}

func (suite RouteSuite) TestRegisterRoutes() {
	app := echo.New()
	healthHandler := &mockHealthCheckHandler{}
	healthHandler.On("Check").Return(nil)
	RegisterRoutes(app, &providers.Handlers{
		Healthcheck: healthHandler,
		Config:      &mockConfigHandler{},
		Network:     &mockNetworkHandler{},
		Status:      &mockStatusHandler{},
	})

	e := []string{}
	for _, r := range app.Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	suite.Contains(e, "GET /api/network")
	suite.Contains(e, "POST /api/network")
	suite.Contains(e, "GET /api/healthz")
	suite.Contains(e, "GET /api/appconfig")
	suite.Contains(e, "GET /api/status/:name")
}

func TestRouteSuite(t *testing.T) {
	suite.Run(t, new(RouteSuite))
}
