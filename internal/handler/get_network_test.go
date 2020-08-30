package handler

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetNetwork(t *testing.T) {
	app := echo.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)
	h := newMockHandler()
	if assert.NoError(t, h.getNetwork(ctx)) {
		dec := json.NewDecoder(strings.NewReader(rec.Body.String()))
		var net models.Network
		if assert.NoError(t, dec.Decode(&net)) {
			assert.Equal(t, "test-network", net.Network, "Get Network should return 'test-network'")
			assert.Equal(t, "a", net.Sites[0].FriendlyName, "The first site in the list should have FriendlyName 'a'")
			assert.Equal(t, uint(2), net.Sites[0].ID, "The first site in the list should be ID '2'")
			assert.Equal(t, "z", net.Sites[1].FriendlyName, "The second site in the list should have FriendlyName 'z'")
			assert.Equal(t, uint(1), net.Sites[1].ID, "The second site in the list should have ID '1'")
		}

	}
}

func TestGetNetworkError(t *testing.T) {
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
			Store: mockStore{
				CreateNetworkFunc: func(net *models.Network) error { return errors.New("not implemented") },
				GetNetworkFunc:    func(net *models.Network) error { return c.Err },
				GetSiteFunc:       func(site *models.Site) error { return errors.New("not implemented") },
			},
		}
		err := h.getNetwork(ctx)
		assert.EqualError(t, err, c.Expected.Error())
	}
}

func TestAddGetNetworkRoute(t *testing.T) {
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	h := newMockHandler()
	h.AddGetNetworkRoute(app)
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	assert.Contains(t, e, "GET /api/network")
}
