package status

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/jarcoal/httpmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockSiteStore struct {
	GetFunc func(*model.Site)
}

func (m *mockSiteStore) GetSite(site *model.Site) {
	m.GetFunc(site)
}

func TestHttp(t *testing.T) {
	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))

	url, err := url.Parse("https://my.test.site")
	assert.NoError(t, err)
	err = testHTTP(url)
	assert.NoError(t, err, "https://my.test.site should not error")
}

func TestHttpTimeout(t *testing.T) {
	os.Setenv("GSS_TIMEOUT", "100")
	config.InitConfig("", "no-file")
	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://timeout.test.site", func(req *http.Request) (*http.Response, error) {
		time.Sleep(2 * time.Second)
		return httpmock.NewStringResponse(200, "success"), nil
	})

	url, err := url.Parse("https://timeout.test.site")
	assert.NoError(t, err)
	err = testHTTP(url)
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
		Site     model.Site
		IsUp     bool
		HasError bool
	}{
		{Site: model.Site{URI: "https://my.test.site"}, IsUp: true, HasError: false},
		{Site: model.Site{URI: "https://my.fail.site"}, IsUp: false, HasError: true},
		{Site: model.Site{URI: "https://^^invalidurl^^"}, IsUp: false, HasError: true},
		{Site: model.Site{URI: "ssh://localhost:12345"}, IsUp: true, HasError: false},
		{Site: model.Site{URI: "ssh://localhost:1234"}, IsUp: false, HasError: true},
		{Site: model.Site{URI: "https://err.test.site"}, IsUp: false, HasError: true},
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
	var store mockSiteStore
	h := Handler{Store: &store}

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
		{ID: "1", URI: "https://my.fail.site", IsUp: false, Error: nil},
		{ID: "1", URI: "https://^^invalidurl^^", IsUp: false, Error: nil},
		{ID: "1", URI: "ssh://localhost:12345", IsUp: true, Error: nil},
		{ID: "1", URI: "ssh://localhost:1234", IsUp: false, Error: nil},
		{ID: "1", URI: "https://500.test.site", IsUp: false, Error: nil},
		{ID: "abc", URI: "https://400.test.site", IsUp: false, Error: echo.ErrBadRequest},
		{ID: "", URI: "https://no-id.test.site", IsUp: false, Error: echo.ErrBadRequest},
	}

	for _, c := range cases {
		store.GetFunc = func(site *model.Site) {
			site.ID = 1
			site.URI = c.URI
		}

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		ctx := app.NewContext(req, rec)
		ctx.SetPath("/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(c.ID)
		err := h.Get(ctx)
		if c.Error != nil {
			assert.EqualError(t, err, c.Error.Error(), "%s should return %s", c.URI, c.Error.Error())
		}
		if c.IsUp {
			assert.Contains(t, rec.Body.String(), `"isUp":true`, "%s should be up", c.URI)
		}
	}
}

func TestGetStatus(t *testing.T) {
	var store mockSiteStore
	handler := Handler{Store: &store}

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
		store.GetFunc = func(site *model.Site) {
			site.ID = 1
			site.URI = c.URI
		}
		res, err := getStatus(handler, 1)
		if c.Error {
			assert.Error(t, err, "%s should error", c.URI)
		} else {
			assert.NoError(t, err, "%s should not error", c.URI)
		}
		assert.Equal(t, c.IsUp, res.IsUp, "%s should have isUp=%t", c.URI, c.IsUp)
	}
}
