package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/eiladin/go-simple-startpage/pkg/usecases/network"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockNetworkUseCase struct {
	mock.Mock
}

func (m *mockNetworkUseCase) Get() (*models.Network, error) {
	args := m.Called()
	data := args.Get(0).(models.Network)
	return &data, args.Error(1)
}

func (m *mockNetworkUseCase) Create(net *models.Network) error {
	args := m.Called(net)
	return args.Error(0)
}

type NetworkSuite struct {
	suite.Suite
}

func (suite *NetworkSuite) TestCreate() {
	app := echo.New()
	body := `{ "network": "test-network" }`

	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	uc := new(mockNetworkUseCase)
	uc.On("Create", &models.Network{Network: "test-network"}).Return(nil)
	h := NetworkHandler{NetworkUseCase: uc}

	if suite.NoError(h.Create(ctx)) {
		suite.Equal(http.StatusCreated, rec.Code, "Create should return a 201")
		uc.AssertExpectations(suite.T())
	}
}

func (suite *NetworkSuite) TestCreateError() {
	cases := []struct {
		body        string
		networkName string
		err         error
	}{
		{body: `{"network":"test"}`, networkName: "test", err: echo.ErrInternalServerError},
		{body: "", networkName: "", err: echo.ErrBadRequest},
	}

	app := echo.New()
	rec := httptest.NewRecorder()

	for _, c := range cases {
		req := httptest.NewRequest("POST", "/", strings.NewReader(c.body))
		req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ctx := app.NewContext(req, rec)
		uc := new(mockNetworkUseCase)
		if !errors.Is(c.err, echo.ErrBadRequest) {
			uc.On("Create", &models.Network{Network: c.networkName}).Return(c.err)
		}

		h := NetworkHandler{NetworkUseCase: uc}
		err := h.Create(ctx)
		suite.Equal(err, c.err)
		uc.AssertExpectations(suite.T())
	}
}

func (suite *NetworkSuite) TestGet() {
	app := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)
	uc := new(mockNetworkUseCase)
	uc.On("Get").Return(
		models.Network{
			Network: "test-network",
			Sites: []models.Site{
				{ID: 1, FriendlyName: "z"},
				{ID: 2, FriendlyName: "a"},
			},
		}, nil)

	ns := NetworkHandler{NetworkUseCase: uc}

	if suite.NoError(ns.Get(ctx)) {
		uc.AssertExpectations(suite.T())
		dec := json.NewDecoder(strings.NewReader(rec.Body.String()))
		var net models.Network
		if suite.NoError(dec.Decode(&net)) {
			suite.Equal("test-network", net.Network, "Get Network should return 'test-network'")
			suite.Len(net.Sites, 2, "There should be 2 sites")
			suite.Equal("z", net.Sites[0].FriendlyName, "The first site in the list should have FriendlyName 'z'")
			suite.Equal(uint(1), net.Sites[0].ID, "The first site in the list should be ID '1'")
			suite.Equal("a", net.Sites[1].FriendlyName, "The second site in the list should have FriendlyName 'a'")
			suite.Equal(uint(2), net.Sites[1].ID, "The second site in the list should have ID '2'")
		}

	}
}

func (suite *NetworkSuite) TestGetError() {
	cases := []struct {
		Err      error
		Expected error
	}{
		{Err: errors.New("not implemented"), Expected: echo.ErrInternalServerError},
		{Err: network.ErrNotFound, Expected: echo.ErrNotFound},
	}
	app := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	for _, c := range cases {
		uc := new(mockNetworkUseCase)
		uc.On("Create", mock.Anything).Return(errors.New("not implemented"))
		uc.On("Get").Return(models.Network{}, c.Err)
		ns := NetworkHandler{NetworkUseCase: uc}
		err := ns.Get(ctx)
		suite.EqualError(err, c.Expected.Error())
	}
}

func TestNetworkSuite(t *testing.T) {
	suite.Run(t, new(NetworkSuite))
}
