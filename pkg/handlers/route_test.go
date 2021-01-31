package handlers

import (
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/providers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type RouteSuite struct {
	suite.Suite
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
