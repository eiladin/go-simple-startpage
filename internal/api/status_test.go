package api

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/jarcoal/httpmock"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/stretchr/testify/assert"
)

func TestHttp(t *testing.T) {
	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))

	url, err := url.Parse("https://my.test.site")
	assert.NoError(t, err)
	err = testHTTP(0, url)
	assert.NoError(t, err, "https://my.test.site should not error")
}

func TestHttpTimeout(t *testing.T) {
	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://timeout.test.site", func(req *http.Request) (*http.Response, error) {
		time.Sleep(2 * time.Second)
		return httpmock.NewStringResponse(200, "success"), nil
	})

	url, err := url.Parse("https://timeout.test.site")
	assert.NoError(t, err)
	err = testHTTP(100, url)
	assert.Error(t, err, "https://timeout.test.site should timeout")
	os.Unsetenv("GSS_TIMEOUT")
}

func TestTCP(t *testing.T) {
	ln, err := net.Listen("tcp", "[::]:1234")
	assert.NoError(t, err)
	defer ln.Close()

	url, err := url.Parse("ssh://localhost:1234")
	assert.NoError(t, err)
	err = testSSH(url)
	assert.NoError(t, err, "ssh://localhost:1234 should not error")
}

func TestGetIP(t *testing.T) {
	url, err := url.Parse("http://localhost")
	assert.NoError(t, err)
	ip := getIP(url)
	assert.Contains(t, []string{"127.0.0.1", "::1"}, ip, "http://localhost should return the following ips: [127.0.0.1, ::1]")
}

func TestUpdateStatus(t *testing.T) {
	cases := []struct {
		Site     models.Site
		IsUp     bool
		HasError bool
	}{
		{Site: models.Site{URI: "https://my.test.site"}, IsUp: true, HasError: false},
		{Site: models.Site{URI: "https://my.fail.site"}, IsUp: false, HasError: true},
		{Site: models.Site{URI: "https://^^invalidurl^^"}, IsUp: false, HasError: true},
		{Site: models.Site{URI: "ssh://localhost:12345"}, IsUp: true, HasError: false},
		{Site: models.Site{URI: "ssh://localhost:1234"}, IsUp: false, HasError: true},
		{Site: models.Site{URI: "https://err.test.site"}, IsUp: false, HasError: true},
	}

	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))
	httpmock.RegisterResponder("GET", "https://err.test.site", httpmock.NewStringResponder(101, "redirect"))

	ln, err := net.Listen("tcp", "[::]:12345")
	assert.NoError(t, err)
	defer ln.Close()

	for _, c := range cases {
		err := updateStatus(0, &c.Site)
		if c.HasError {
			assert.Error(t, err, "site: %s should error", c.Site.URI)
			assert.False(t, c.IsUp, "site: %s should not be up", c.Site.URI)
		} else {
			assert.NoError(t, err, "site: %s should not error", c.Site.URI)
			assert.Equal(t, c.IsUp, c.Site.IsUp, "site: %s should be up", c.Site.URI)
		}
	}
}

func TestGetStatus(t *testing.T) {
	app := echo.New()
	var s mockStore
	h := handler{Store: &s, Config: &models.Config{}}

	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))
	httpmock.RegisterResponder("GET", "https://err.test.site", httpmock.NewStringResponder(101, "fail"))
	httpmock.RegisterResponder("GET", "https://bigid.test.site", httpmock.NewStringResponder(200, "success"))

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
		{ID: "1", URI: "https://my.fail.site", IsUp: false, Error: nil},
		{ID: "1", URI: "https://^^invalidurl^^", IsUp: false, Error: nil},
		{ID: "1", URI: "ssh://localhost:12345", IsUp: true, Error: nil},
		{ID: "1", URI: "ssh://localhost:1234", IsUp: false, Error: nil},
		{ID: "1", URI: "https://500.test.site", IsUp: false, Error: nil},
		{ID: "abc", URI: "https://400.test.site", IsUp: false, Error: echo.ErrBadRequest},
		{ID: "", URI: "https://no-id.test.site", IsUp: false, Error: echo.ErrBadRequest},
		{ID: "12345", URI: "https://bigid.test.site", IsUp: false, Error: echo.ErrNotFound},
		{ID: "0", URI: "https://my.test.site", IsUp: false, Error: echo.ErrBadRequest},
	}

	for _, c := range cases {
		s.GetSiteFunc = func(site *models.Site) error {
			if site.ID != 1 {
				return store.ErrNotFound
			}
			site.ID = 1
			site.URI = c.URI
			return nil
		}

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		ctx := app.NewContext(req, rec)
		ctx.SetPath("/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(c.ID)
		err := h.getStatus(ctx)
		if c.Error != nil {
			assert.EqualError(t, err, c.Error.Error(), "%s should return %s", c.URI, c.Error.Error())
		}
		if c.IsUp {
			assert.Contains(t, rec.Body.String(), `"isUp":true`, "%s should be up", c.URI)
		}
	}
}

func TestAddStatusRoutes(t *testing.T) {
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	h := handler{ApiRoot: app, Store: &mockStore{}}
	h.addStatusRoutes()
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	assert.Contains(t, e, "GET /api/status/:id")
}