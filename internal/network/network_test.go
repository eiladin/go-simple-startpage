package network

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockNetworkService struct {
	CreateNetworkFunc func(*model.Network)
	FindNetworkFunc   func(*model.Network)
}

func (m *mockNetworkService) CreateNetwork(net *model.Network) {
	m.CreateNetworkFunc(net)
}

func (m *mockNetworkService) FindNetwork(net *model.Network) {
	m.FindNetworkFunc(net)
}

func getMockHandler() Handler {
	store := mockNetworkService{
		CreateNetworkFunc: func(net *model.Network) {
			net.ID = 12345
		},
		FindNetworkFunc: func(net *model.Network) {
			net.ID = 12345
			net.Network = "test-network"
		},
	}
	return Handler{NetworkService: &store}
}

func TestCreateHandler(t *testing.T) {
	app := echo.New()
	body := `{ "network": "test-network" }`

	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	h := getMockHandler()
	if assert.NoError(t, h.Create(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "{\"id\":12345}\n", rec.Body.String())
	}
}

func TestCreateError(t *testing.T) {
	app := echo.New()
	req := httptest.NewRequest("POST", "/", strings.NewReader(``))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)
	h := getMockHandler()
	err := h.Create(ctx)
	assert.Error(t, err)
	assert.EqualError(t, err, echo.ErrBadRequest.Error())
}

func TestGet(t *testing.T) {
	app := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)
	h := getMockHandler()
	if assert.NoError(t, h.Get(ctx)) {
		assert.Equal(t, "{\"network\":\"test-network\",\"links\":null,\"sites\":null}\n", rec.Body.String())
	}
}
