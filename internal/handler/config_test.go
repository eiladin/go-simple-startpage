package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetAppConfig(t *testing.T) {
	viper.Reset()
	app := echo.New()
	c := config.New("1.2.3", "not-found")
	h := Config{Store: c}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	if assert.NoError(t, h.Get(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code, "Status code should be 200")
		assert.Equal(t, "{\"version\":\"1.2.3\"}\n", rec.Body.String(), "Version should be 1.2.3")
	}
}

func TestConfigHandler(t *testing.T) {
	viper.Reset()
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	h := Config{Store: models.Config{
		Version: "1.2.3",
	}}
	h.Register(app)
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Path)
	}
	assert.Contains(t, e, "/api/appconfig")
}
