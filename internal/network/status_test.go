package network

import (
	"fmt"
	"net"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
	"github.com/gofiber/fiber"
	"github.com/jarcoal/httpmock"
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
	var store mockNetworkService
	handler := Handler{NetworkService: &store}

	app := fiber.New()
	app.Get("/:id", handler.GetStatus)

	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))
	httpmock.RegisterResponder("GET", "https://err.test.site", httpmock.NewStringResponder(101, "fail"))

	ln, err := net.Listen("tcp", "[::]:12345")
	assert.NoError(t, err)
	defer ln.Close()

	cases := []struct {
		Desc       string
		ID         string
		URI        string
		IsUp       bool
		StatusCode int
	}{
		{Desc: "should be up", ID: "1", URI: "https://my.test.site", IsUp: true, StatusCode: 200},
		{Desc: "should not be up", ID: "1", URI: "https://my.fail.site", IsUp: false, StatusCode: 500},
		{Desc: "should not be up", ID: "1", URI: "https://^^invalidurl^^", IsUp: false, StatusCode: 500},
		{Desc: "should be up", ID: "1", URI: "ssh://localhost:12345", IsUp: true, StatusCode: 200},
		{Desc: "should not be up", ID: "1", URI: "ssh://localhost:1234", IsUp: false, StatusCode: 500},
		{Desc: "should not be up", ID: "1", URI: "https://err.test.site", IsUp: false, StatusCode: 500},
		{Desc: "should not be up", ID: "abc", URI: "https://err.test.site", IsUp: false, StatusCode: 500},
		{Desc: "should not be up", ID: "", URI: "https://err.test.site", IsUp: false, StatusCode: 404},
	}

	for _, c := range cases {

		store.FindSiteFunc = func(site *interfaces.Site) {
			site.ID = 1
			site.URI = c.URI
		}

		req := httptest.NewRequest("GET", fmt.Sprintf("/%s", c.ID), nil)
		resp, err := app.Test(req)
		assert.NoError(t, err, fmt.Sprintf("%s %s", c.URI, c.Desc))
		assert.Equal(t, c.StatusCode, resp.StatusCode, fmt.Sprintf("%s %s", c.URI, c.Desc))
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
		Desc  string
		URI   string
		IsUp  bool
		Error bool
	}{
		{Desc: "should be up", URI: "https://my.test.site", IsUp: true, Error: false},
		{Desc: "should not be up", URI: "https://my.fail.site", IsUp: false, Error: true},
		{Desc: "should not be up", URI: "https://^^invalidurl^^", IsUp: false, Error: true},
		{Desc: "should be up", URI: "ssh://localhost:12345", IsUp: true, Error: false},
		{Desc: "should not be up", URI: "ssh://localhost:1234", IsUp: false, Error: true},
		{Desc: "should not be up", URI: "https://err.test.site", IsUp: false, Error: true},
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
		assert.Equal(t, c.IsUp, res.IsUp, fmt.Sprintf("%s %s", c.URI, c.Desc))
	}
}
