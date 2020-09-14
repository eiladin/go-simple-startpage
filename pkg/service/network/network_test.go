package network

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
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

	ns := NewNetworkService(&model.Config{}, &mockStore{
		CreateNetworkFunc: func(net *model.Network) error {
			net.ID = 12345
			return nil
		},
	})
	if suite.NoError(ns.Create(ctx)) {
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
		ns := NewNetworkService(&model.Config{}, &mockStore{
			CreateNetworkFunc: func(net *model.Network) error { return errors.New("not implemented") },
		})
		err := ns.Create(ctx)
		suite.EqualError(err, c.Err.Error())
	}
}

func (suite *NetworkServiceSuite) TestGet() {
	app := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)
	ns := NewNetworkService(&model.Config{}, &mockStore{
		GetNetworkFunc: func(net *model.Network) error {
			net.Network = "test-network"
			net.Sites = []model.Site{
				{ID: 1, FriendlyName: "z"},
				{ID: 2, FriendlyName: "a"},
			}
			return nil
		},
	})
	if suite.NoError(ns.Get(ctx)) {
		dec := json.NewDecoder(strings.NewReader(rec.Body.String()))
		var net model.Network
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
		ns := NewNetworkService(&model.Config{}, &mockStore{
			CreateNetworkFunc: func(net *model.Network) error { return errors.New("not implemented") },
			GetNetworkFunc:    func(net *model.Network) error { return c.Err },
			GetSiteFunc:       func(site *model.Site) error { return errors.New("not implemented") },
		})
		err := ns.Get(ctx)
		suite.EqualError(err, c.Expected.Error())
	}
}

func (suite *NetworkServiceSuite) TestRegister() {
	app := echo.New()
	NewNetworkService(&model.Config{}, &mockStore{}).Register(app)
	e := []string{}
	for _, r := range app.Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	suite.Contains(e, "GET /api/network")
	suite.Contains(e, "POST /api/network")
}

func TestNetworkServiceSuite(t *testing.T) {
	suite.Run(t, new(NetworkServiceSuite))
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
