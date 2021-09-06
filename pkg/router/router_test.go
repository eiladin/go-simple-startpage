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

type mockHealthcheckUseCase struct {
	mock.Mock
}

func (m *mockHealthcheckUseCase) Check() http.Handler {
	m.Called()
	return nil
}

type mockNetworkUseCase struct {
	mock.Mock
}

func (m *mockNetworkUseCase) Get() (*network.Network, error) {
	args := m.Called()
	data := args.Get(0).(network.Network)
	return &data, args.Error(1)
}

func (m *mockNetworkUseCase) Create(net *network.Network) error {
	args := m.Called(net)
	return args.Error(0)
}

type mockConfigUseCase struct {
	mock.Mock
}

func (m *mockConfigUseCase) Get() (*config.Config, error) {
	args := m.Called()
	data := args.Get(0).(config.Config)
	return &data, args.Error(1)
}

type HandlerSuite struct {
	suite.Suite
}

type mockStatusUseCase struct {
	mock.Mock
}

func (m *mockStatusUseCase) Get(name string) (*status.Status, error) {
	args := m.Called(name)
	data := args.Get(0).(status.Status)
	return &data, args.Error(1)
}

func (suite RouteSuite) TestRegisterRoutes() {
	app := echo.New()
	huc := &mockHealthcheckUseCase{}
	huc.On("Check").Return(nil)
	RegisterRoutes(app, &providers.Provider{
		Healthcheck: huc,
		Config:      &mockConfigUseCase{},
		Network:     &mockNetworkUseCase{},
		Status:      &mockStatusUseCase{},
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
