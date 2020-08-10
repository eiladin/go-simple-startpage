package network

import (
	"fmt"
	"net"
	"net/url"
	"testing"

	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
	"github.com/gofiber/fiber"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
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

func TestUpdateStatusHandler(t *testing.T) {
	app := fiber.New()
	var store mockNetworkService
	handler := Handler{NetworkService: &store}

	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))
	httpmock.RegisterResponder("GET", "https://err.test.site", httpmock.NewStringResponder(101, "redirect"))

	ln, err := net.Listen("tcp", "[::]:12345")
	assert.NoError(t, err)
	defer ln.Close()

	cases := []struct {
		Desc       string
		Body       string
		IsUp       bool
		StatusCode int
	}{
		{Desc: "https://my.test.site should be up", Body: `{"uri":"https://my.test.site"}`, IsUp: true, StatusCode: 200},
		{Desc: "https://my.fail.site should not be up", Body: `{"uri":"https://my.fail.site"}`, IsUp: false, StatusCode: 500},
		{Desc: "https://^^invalidurl^^ should not be up", Body: `{"uri":"https://^^invalidurl^^"}`, IsUp: false, StatusCode: 500},
		{Desc: "ssh://localhost:12345 should be up", Body: `{"uri":"ssh://localhost:12345"}`, IsUp: true, StatusCode: 200},
		{Desc: "ssh://localhost:1234 should not be up", Body: `{"uri":"ssh://localhost:1234"}`, IsUp: false, StatusCode: 500},
		{Desc: "https://err.test.site should not be up", Body: `{"uri":"https://err.test.site"}`, IsUp: false, StatusCode: 500},
		{Desc: "invalid json should return a 400", Body: `{invalid json}`, IsUp: false, StatusCode: 400},
	}

	for _, c := range cases {
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		ctx.Fasthttp.Request.Header.SetContentType(fiber.MIMEApplicationJSON)
		ctx.Fasthttp.Request.SetBody([]byte(c.Body))
		ctx.Fasthttp.Request.Header.SetContentLength(len(c.Body))
		handler.UpdateStatus(ctx)

		assert.Equal(t, c.StatusCode, ctx.Fasthttp.Response.StatusCode())
		switch c.StatusCode {
		case 500:
			assert.Contains(t, string(ctx.Fasthttp.Response.Body()), `"isUp":false`, c.Desc)
		case 200:
			assert.Contains(t, string(ctx.Fasthttp.Response.Body()), `"isUp":true`, c.Desc)
		}
		app.ReleaseCtx(ctx)
	}
}
