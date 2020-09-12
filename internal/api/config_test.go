package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/stretchr/testify/suite"
)

type ConfigServiceSuite struct {
	suite.Suite
}

func (suite *ConfigServiceSuite) TestGet() {
	app := echo.New()
	c := models.Config{Version: "1.2.3"}
	cs := NewConfigService(&c)
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	ctx := app.NewContext(req, rec)

	if suite.NoError(cs.Get(ctx)) {
		suite.Equal(http.StatusOK, rec.Code, "Status code should be 200")
		suite.Equal("{\"version\":\"1.2.3\"}\n", rec.Body.String(), "Version should be 1.2.3")
	}
}

func (suite *ConfigServiceSuite) TestRegister() {
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	c := models.Config{Version: "1.2.3"}
	NewConfigService(&c).Register(app)
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Path)
	}
	suite.Contains(e, "/api/appconfig")
}

func TestConfigServiceSuite(t *testing.T) {
	suite.Run(t, new(ConfigServiceSuite))
}
