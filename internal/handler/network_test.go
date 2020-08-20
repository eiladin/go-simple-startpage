package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func getMockNetwork() Network {
	store := mockStore{
		CreateNetworkFunc: func(net *model.Network) {
			net.ID = 12345
		},
		GetNetworkFunc: func(net *model.Network) {
			net.ID = 12345
			net.Network = "test-network"
			net.Sites = []model.Site{
				{ID: 1, FriendlyName: "z"},
				{ID: 2, FriendlyName: "a"},
			}
		},
		GetSiteFunc: func(site *model.Site) {},
	}
	return Network{Store: &store}
}

func TestCreateHandler(t *testing.T) {
	app := echo.New()
	body := `{ "network": "test-network" }`

	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	h := getMockNetwork()
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
	h := getMockNetwork()
	err := h.Create(ctx)
	assert.EqualError(t, err, echo.ErrBadRequest.Error(), "Invalid Post should return a BadRequest (400)")
}

func TestGet(t *testing.T) {
	app := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)
	h := getMockNetwork()
	if assert.NoError(t, h.Get(ctx)) {
		dec := json.NewDecoder(strings.NewReader(rec.Body.String()))
		var net model.Network
		if assert.NoError(t, dec.Decode(&net)) {
			assert.Equal(t, "test-network", net.Network, "Get Network should return 'test-network'")
			assert.Equal(t, "a", net.Sites[0].FriendlyName, "The first site in the list should have FriendlyName 'a'")
			assert.Equal(t, uint(2), net.Sites[0].ID, "The first site in the list should be ID '2'")
			assert.Equal(t, "z", net.Sites[1].FriendlyName, "The second site in the list should have FriendlyName 'z'")
			assert.Equal(t, uint(1), net.Sites[1].ID, "The second site in the list should have ID '1'")
		}

	}
}
