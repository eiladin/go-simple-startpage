package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetAppConfig(t *testing.T) {
	viper.Reset()
	app := echo.New()
	config.InitConfig("1.2.3", "not-found")
	h := Config{}
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	c := app.NewContext(req, rec)

	if assert.NoError(t, h.Get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code, "Status code should be 200")
		assert.Equal(t, "{\"version\":\"1.2.3\"}\n", rec.Body.String(), "Version should be 1.2.3")
	}
}
