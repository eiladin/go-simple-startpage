package handlers

import (
	"encoding/json"
	"errors"
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
	cases := []struct {
		body        string
		networkName string
		err         error
	}{
		{body: `{"network":"test"}`, networkName: "test", err: nil},
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
	cases := []struct {
		network  models.Network
		err      error
		expected error
	}{
		{
			network: models.Network{
				Network: "test-network",
				Sites: []models.Site{
					{ID: 1, Name: "z"},
					{ID: 2, Name: "a"},
				},
			},
			err:      nil,
			expected: nil,
		},
		{
			network:  models.Network{},
			err:      errors.New("not implemented"),
			expected: echo.ErrInternalServerError,
		},
		{
			network:  models.Network{},
			err:      network.ErrNotFound,
			expected: echo.ErrNotFound,
		},
	}

	app := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	for _, c := range cases {
		uc := new(mockNetworkUseCase)
		uc.On("Get").Return(c.network, c.err)
		ns := NetworkHandler{NetworkUseCase: uc}

		err := ns.Get(ctx)

		uc.AssertExpectations(suite.T())
		if c.expected != nil {
			suite.EqualError(err, c.expected.Error())
		} else {
			dec := json.NewDecoder(strings.NewReader(rec.Body.String()))
			var net models.Network
			if suite.NoError(dec.Decode(&net)) {
				suite.Equal("test-network", net.Network, "Get Network should return 'test-network'")
				suite.Len(net.Sites, 2, "There should be 2 sites")
				suite.Equal("z", net.Sites[0].Name, "The first site in the list should have Name 'z'")
				suite.Equal("a", net.Sites[1].Name, "The second site in the list should have Name 'a'")
			}
		}
	}
}

func TestNetworkSuite(t *testing.T) {
	suite.Run(t, new(NetworkSuite))
}
