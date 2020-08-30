package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	app := echo.New()
	c := models.Config{Version: "1.2.3"}
	h := handler{Config: &c}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	if assert.NoError(t, h.getConfig(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code, "Status code should be 200")
		assert.Equal(t, "{\"version\":\"1.2.3\"}\n", rec.Body.String(), "Version should be 1.2.3")
	}
}

func TestAddGetConfigRoute(t *testing.T) {
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	c := models.Config{Version: "1.2.3"}
	h := handler{Config: &c}
	h.addGetConfigRoute(app)
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Path)
	}
	assert.Contains(t, e, "/api/appconfig")
}
