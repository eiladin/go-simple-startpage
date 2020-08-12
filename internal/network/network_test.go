package network

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type mockNetworkService struct {
	CreateNetworkFunc func(*interfaces.Network)
	FindNetworkFunc   func(*interfaces.Network)
	FindSiteFunc      func(*interfaces.Site)
}

func (m *mockNetworkService) CreateNetwork(net *interfaces.Network) {
	m.CreateNetworkFunc(net)
}

func (m *mockNetworkService) FindNetwork(net *interfaces.Network) {
	m.FindNetworkFunc(net)
}

func (m *mockNetworkService) FindSite(site *interfaces.Site) {
	m.FindSiteFunc(site)
}

func getMockHandler() Handler {
	store := mockNetworkService{
		CreateNetworkFunc: func(net *interfaces.Network) {
			net.ID = 12345
		},
		FindNetworkFunc: func(net *interfaces.Network) {
			net.ID = 12345
			net.Network = "test-network"
		},
	}
	return Handler{NetworkService: &store}
}

func TestNewNetworkHandler(t *testing.T) {
	app := echo.New()
	body := `{ "network": "test-network" }`

	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	h := getMockHandler()
	if assert.NoError(t, h.NewNetwork(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "{\"id\":\"12345\"}", rec.Body.String())
	}
}

func TestNewNetworkError(t *testing.T) {
	app := echo.New()
	req := httptest.NewRequest("POST", "/", strings.NewReader(``))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)
	h := getMockHandler()
	err := h.NewNetwork(ctx)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, echo.ErrBadRequest))
}

func TestGetNetwork(t *testing.T) {
	app := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)
	h := getMockHandler()
	if assert.NoError(t, h.GetNetwork(ctx)) {
		assert.Equal(t, "{\"network\":\"test-network\",\"links\":null,\"sites\":null}\n", rec.Body.String())
	}
}
