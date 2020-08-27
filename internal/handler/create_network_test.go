package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateHandler(t *testing.T) {
	app := echo.New()
	body := `{ "network": "test-network" }`

	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	h := getMockHandler()
	if assert.NoError(t, h.CreateNetwork(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code, "Create should return a 201")
		assert.Equal(t, "{\"id\":12345}\n", rec.Body.String(), "Create should return an ID")
	}
}

func TestCreateError(t *testing.T) {
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
			Store: mockStore{
				CreateNetworkFunc: func(net *models.Network) error { return errors.New("not implemented") },
				GetNetworkFunc:    func(net *models.Network) error { return errors.New("not implemented") },
				GetSiteFunc:       func(site *models.Site) error { return errors.New("not implemented") },
			},
		}
		err := h.CreateNetwork(ctx)
		assert.EqualError(t, err, c.Err.Error())
	}
}

func TestCreateNetworkRegister(t *testing.T) {
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	h := getMockHandler()
	h.AddPostNetworkRoute(app)
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	assert.Contains(t, e, "POST /api/network")
}
