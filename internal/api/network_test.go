package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/stretchr/testify/suite"
)

type NetworkServiceSuite struct {
	suite.Suite
}

func (suite *NetworkServiceSuite) TestCreate() {
	app := echo.New()
	body := `{ "network": "test-network" }`

	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	h := handler{Store: &mockStore{
		CreateNetworkFunc: func(net *models.Network) error {
			net.ID = 12345
			return nil
		},
	}}
	if suite.NoError(h.createNetwork(ctx)) {
		suite.Equal(http.StatusCreated, rec.Code, "Create should return a 201")
		suite.Equal("{\"id\":12345}\n", rec.Body.String(), "Create should return an ID")
	}
}

func (suite *NetworkServiceSuite) TestCreateError() {
	cases := []struct {
		Body string
		Err  error
	}{
		{Body: `{"network":"test"}`, Err: echo.ErrInternalServerError},
		{Body: "", Err: echo.ErrBadRequest},
	}

	app := echo.New()
	rec := httptest.NewRecorder()

	for _, c := range cases {
		req := httptest.NewRequest("POST", "/", strings.NewReader(c.Body))
		req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := app.NewContext(req, rec)
		h := handler{
			Store: &mockStore{
				CreateNetworkFunc: func(net *models.Network) error { return errors.New("not implemented") },
			},
		}
		err := h.createNetwork(ctx)
		suite.EqualError(err, c.Err.Error())
	}
}

func (suite *NetworkServiceSuite) TestGet() {
	app := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)
	h := handler{Store: &mockStore{
		GetNetworkFunc: func(net *models.Network) error {
			net.Network = "test-network"
			net.Sites = []models.Site{
				{ID: 1, FriendlyName: "z"},
				{ID: 2, FriendlyName: "a"},
			}
			return nil
		},
	}}
	if suite.NoError(h.getNetwork(ctx)) {
		dec := json.NewDecoder(strings.NewReader(rec.Body.String()))
		var net models.Network
		if suite.NoError(dec.Decode(&net)) {
			suite.Equal("test-network", net.Network, "Get Network should return 'test-network'")
			suite.Len(net.Sites, 2, "There should be 2 sites")
			suite.Equal("a", net.Sites[0].FriendlyName, "The first site in the list should have FriendlyName 'a'")
			suite.Equal(uint(2), net.Sites[0].ID, "The first site in the list should be ID '2'")
			suite.Equal("z", net.Sites[1].FriendlyName, "The second site in the list should have FriendlyName 'z'")
			suite.Equal(uint(1), net.Sites[1].ID, "The second site in the list should have ID '1'")
		}

	}
}

func (suite *NetworkServiceSuite) TestGetError() {
	cases := []struct {
		Err      error
		Expected error
	}{
		{Err: errors.New("not implemented"), Expected: echo.ErrInternalServerError},
		{Err: store.ErrNotFound, Expected: echo.ErrNotFound},
	}
	app := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	for _, c := range cases {
		h := handler{
			Store: &mockStore{
				CreateNetworkFunc: func(net *models.Network) error { return errors.New("not implemented") },
				GetNetworkFunc:    func(net *models.Network) error { return c.Err },
				GetSiteFunc:       func(site *models.Site) error { return errors.New("not implemented") },
			},
		}
		err := h.getNetwork(ctx)
		suite.EqualError(err, c.Expected.Error())
	}
}

func (suite *NetworkServiceSuite) TestRegister() {
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	h := handler{Store: &mockStore{}}
	h.ApiRoot = app
	h.addNetworkRoutes()
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	suite.Contains(e, "GET /api/network")
	suite.Contains(e, "POST /api/network")
}

func TestNetworkServiceSuite(t *testing.T) {
	suite.Run(t, new(NetworkServiceSuite))
}
