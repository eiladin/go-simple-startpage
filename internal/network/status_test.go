package network

import (
	"fmt"
	"net"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
	"github.com/jarcoal/httpmock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestHttp(t *testing.T) {
	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))

	site := interfaces.Site{}
	url, err := url.Parse("https://my.test.site")
	assert.NoError(t, err)
	err = testHTTP(&site, url)
	assert.NoError(t, err)
	assert.Equal(t, true, site.IsUp)
}

func TestTCP(t *testing.T) {
	ln, err := net.Listen("tcp", "[::]:1234")
	assert.NoError(t, err)
	defer ln.Close()

	site := interfaces.Site{}
	url, err := url.Parse("ssh://localhost:1234")
	assert.NoError(t, err)
	err = testSSH(&site, url)
	assert.NoError(t, err)
	assert.Equal(t, true, site.IsUp)
}

func TestGetIP(t *testing.T) {
	url, err := url.Parse("http://localhost")
	assert.NoError(t, err)
	ip := getIP(url)
	assert.Contains(t, []string{"127.0.0.1", "::1"}, ip)
}

func TestUpdateStatus(t *testing.T) {
	cases := []struct {
		Site     interfaces.Site
		IsUp     bool
		HasError bool
	}{
		{Site: interfaces.Site{URI: "https://my.test.site"}, IsUp: true, HasError: false},
		{Site: interfaces.Site{URI: "https://my.fail.site"}, IsUp: false, HasError: true},
		{Site: interfaces.Site{URI: "https://^^invalidurl^^"}, IsUp: false, HasError: true},
		{Site: interfaces.Site{URI: "ssh://localhost:12345"}, IsUp: true, HasError: false},
		{Site: interfaces.Site{URI: "ssh://localhost:1234"}, IsUp: false, HasError: true},
		{Site: interfaces.Site{URI: "https://err.test.site"}, IsUp: false, HasError: true},
	}

	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))
	httpmock.RegisterResponder("GET", "https://err.test.site", httpmock.NewStringResponder(101, "redirect"))

	ln, err := net.Listen("tcp", "[::]:12345")
	assert.NoError(t, err)
	defer ln.Close()

	for _, c := range cases {
		err := updateStatus(&c.Site)
		if c.HasError {
			assert.Error(t, err, fmt.Sprintf("site: %s should error", c.Site.URI))
			assert.False(t, c.IsUp, fmt.Sprintf("site: %s should not be up", c.Site.URI))
		} else {
			assert.NoError(t, err, fmt.Sprintf("site: %s should not error", c.Site.URI))
			assert.Equal(t, c.IsUp, c.Site.IsUp, fmt.Sprintf("site: %s should be up", c.Site.URI))
		}
	}
}

func TestGetStatusHandler(t *testing.T) {
	app := echo.New()
	var store mockNetworkService
	h := Handler{NetworkService: &store}

	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))
	httpmock.RegisterResponder("GET", "https://err.test.site", httpmock.NewStringResponder(101, "fail"))
	ln, err := net.Listen("tcp", "[::]:12345")
	assert.NoError(t, err)
	defer ln.Close()

	cases := []struct {
		ID    string
		URI   string
		IsUp  bool
		Error error
	}{
		{ID: "1", URI: "https://my.test.site", IsUp: true, Error: nil},
		{ID: "1", URI: "https://my.fail.site", IsUp: false, Error: echo.ErrInternalServerError},
		{ID: "1", URI: "https://^^invalidurl^^", IsUp: false, Error: echo.ErrInternalServerError},
		{ID: "1", URI: "ssh://localhost:12345", IsUp: true, Error: nil},
		{ID: "1", URI: "ssh://localhost:1234", IsUp: false, Error: echo.ErrInternalServerError},
		{ID: "1", URI: "https://500.test.site", IsUp: false, Error: echo.ErrInternalServerError},
		{ID: "abc", URI: "https://400.test.site", IsUp: false, Error: echo.ErrBadRequest},
		{ID: "", URI: "https://no-id.test.site", IsUp: false, Error: echo.ErrBadRequest},
	}

	for _, c := range cases {
		store.FindSiteFunc = func(site *interfaces.Site) {
			site.ID = 1
			site.URI = c.URI
		}

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		ctx := app.NewContext(req, rec)
		ctx.SetPath("/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(c.ID)
		err := h.GetStatus(ctx)
		if c.Error != nil {
			assert.EqualError(t, err, c.Error.Error(), fmt.Sprintf("%s should return %s", c.URI, c.Error.Error()))
		}
		if c.IsUp {
			assert.Contains(t, rec.Body.String(), `"isUp":true`, fmt.Sprintf("%s should be up", c.URI))
		}
	}
}

func TestGetStatus(t *testing.T) {
	var store mockNetworkService
	handler := Handler{NetworkService: &store}

	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))
	httpmock.RegisterResponder("GET", "https://err.test.site", httpmock.NewStringResponder(101, "fail"))

	ln, err := net.Listen("tcp", "[::]:12345")
	assert.NoError(t, err)
	defer ln.Close()

	cases := []struct {
		URI   string
		IsUp  bool
		Error bool
	}{
		{URI: "https://my.test.site", IsUp: true, Error: false},
		{URI: "https://my.fail.site", IsUp: false, Error: true},
		{URI: "https://^^invalidurl^^", IsUp: false, Error: true},
		{URI: "ssh://localhost:12345", IsUp: true, Error: false},
		{URI: "ssh://localhost:1234", IsUp: false, Error: true},
		{URI: "https://err.test.site", IsUp: false, Error: true},
	}

	for _, c := range cases {
		store.FindSiteFunc = func(site *interfaces.Site) {
			site.ID = 1
			site.URI = c.URI
		}
		res, err := getStatus(handler, 1)
		if c.Error {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, c.IsUp, res.IsUp, fmt.Sprintf("%s should have isUp=%t", c.URI, c.IsUp))
	}
}
