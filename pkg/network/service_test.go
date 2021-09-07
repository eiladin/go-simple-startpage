package network

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockHandler struct {
	mock.Mock
}

func (m *mockHandler) Get() (*Network, error) {
	args := m.Called()
	data := args.Get(0).(Network)
	return &data, args.Error(1)
}

func (m *mockHandler) Create(net *Network) error {
	args := m.Called(net)
	return args.Error(0)
}

type ServiceSuite struct {
	suite.Suite
}

func (suite *ServiceSuite) TestCreate() {
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
		handler := new(mockHandler)
		if !errors.Is(c.err, echo.ErrBadRequest) {
			handler.On("Create", &Network{Network: c.networkName}).Return(c.err)
		}

		h := Service{Handler: handler}
		err := h.Create(ctx)
		suite.Equal(err, c.err)
		handler.AssertExpectations(suite.T())
	}
}

func (suite *ServiceSuite) TestGet() {
	cases := []struct {
		network  Network
		err      error
		expected error
	}{
		{
			network: Network{
				Network: "test-network",
				Sites: []Site{
					{ID: 1, Name: "z"},
					{ID: 2, Name: "a"},
				},
			},
			err:      nil,
			expected: nil,
		},
		{
			network:  Network{},
			err:      errors.New("not implemented"),
			expected: echo.ErrInternalServerError,
		},
		{
			network:  Network{},
			err:      ErrNotFound,
			expected: echo.ErrNotFound,
		},
	}

	app := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	for _, c := range cases {
		handler := new(mockHandler)
		handler.On("Get").Return(c.network, c.err)
		ns := Service{Handler: handler}

		err := ns.Get(ctx)

		handler.AssertExpectations(suite.T())
		if c.expected != nil {
			suite.EqualError(err, c.expected.Error())
		} else {
			dec := json.NewDecoder(strings.NewReader(rec.Body.String()))
			var net Network
			if suite.NoError(dec.Decode(&net)) {
				suite.Equal("test-network", net.Network, "Get Network should return 'test-network'")
				suite.Len(net.Sites, 2, "There should be 2 sites")
				suite.Equal("z", net.Sites[0].Name, "The first site in the list should have Name 'z'")
				suite.Equal("a", net.Sites[1].Name, "The second site in the list should have Name 'a'")
			}
		}
	}
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
