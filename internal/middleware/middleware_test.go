package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/suite"
)

type MiddlewareSuite struct {
	suite.Suite
}

func (suite MiddlewareSuite) TestCSRFSkipper() {
	cases := []struct {
		host     string
		expected bool
	}{
		{host: "localhost", expected: true},
		{host: "myhost", expected: false},
	}

	for _, c := range cases {
		e := echo.New()
		req := httptest.NewRequest("GET", "http://"+c.host, nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		actual := csrfSkipper(ctx)
		suite.Equal(c.expected, actual, "%s should be %t", c.host, c.expected)
	}
}

func (suite MiddlewareSuite) TestGzipSkipper() {
	cases := []struct {
		path     string
		expected bool
	}{
		{path: "/api/status", expected: false},
		{path: "/dashboard", expected: false},
		{path: "/", expected: false},
		{path: "/swagger/index.html", expected: true},
	}

	for _, c := range cases {
		e := echo.New()
		req := httptest.NewRequest("GET", c.path, nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		actual := gzipSkipper(ctx)
		suite.Equal(c.expected, actual, "%s should be %t", c.path, c.expected)
	}
}

func (suite MiddlewareSuite) TestStaticSkipper() {
	cases := []struct {
		path     string
		expected bool
	}{
		{path: "/api/status", expected: true},
		{path: "/dashboard", expected: false},
		{path: "/", expected: false},
		{path: "/swagger/index.html", expected: true},
	}

	for _, c := range cases {
		e := echo.New()
		req := httptest.NewRequest("GET", c.path, nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)
		actual := staticSkipper(ctx)
		suite.Equal(c.expected, actual, "%s should be %t", c.path, c.expected)
	}
}

func (suite MiddlewareSuite) TestLoggerConfig() {
	cases := []struct {
		production bool
		expected   string
	}{
		{production: true, expected: middleware.DefaultLoggerConfig.Format},
		{production: false, expected: "method=${method}, uri=${uri}, status=${status} ${error}\n"}}

	for _, c := range cases {
		cfg := loggerConfig(c.production)
		suite.Equal(c.expected, cfg.Format)
	}
}

func (suite MiddlewareSuite) TestGetMiddleware() {
	m := GetMiddleware(&model.Config{})
	suite.Len(m, 8)
}

func TestMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareSuite))
}
