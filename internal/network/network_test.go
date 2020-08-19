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

type mockNetworkStore struct {
	CreateFunc func(*model.Network)
	GetFunc    func(*model.Network)
}

func (m *mockNetworkStore) CreateNetwork(net *model.Network) {
	m.CreateFunc(net)
}

func (m *mockNetworkStore) GetNetwork(net *model.Network) {
	m.GetFunc(net)
}

func getMockHandler() Handler {
	store := mockNetworkStore{
		CreateFunc: func(net *model.Network) {
			net.ID = 12345
		},
		GetFunc: func(net *model.Network) {
			net.ID = 12345
			net.Network = "test-network"
		},
	}
	return Handler{Store: &store}
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
		assert.Equal(t, http.StatusCreated, rec.Code, "Create should return a 201")
		assert.Equal(t, "{\"id\":12345}\n", rec.Body.String(), "Create should return an ID")
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
	assert.EqualError(t, err, echo.ErrBadRequest.Error(), "Invalid Post should return a BadRequest (400)")
}

func TestGet(t *testing.T) {
	app := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)
	h := getMockHandler()
	if assert.NoError(t, h.Get(ctx)) {
		assert.Equal(t, "{\"network\":\"test-network\",\"links\":null,\"sites\":null}\n", rec.Body.String(), "Get Network should return 'test-network'")
	}
}
