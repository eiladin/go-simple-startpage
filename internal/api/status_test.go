package api

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/models"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/jarcoal/httpmock"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
	"github.com/stretchr/testify/assert"
)

type mockStatusStore struct {
	NewFunc           func() (store.Store, error)
	CreateNetworkFunc func(*models.Network) error
	GetNetworkFunc    func(*models.Network) error
	GetSiteFunc       func(*models.Site) error
}

func (m *mockStatusStore) New() (store.Store, error) {
	return m.NewFunc()
}

func (m *mockStatusStore) CreateNetwork(net *models.Network) error {
	return m.CreateNetworkFunc(net)
}

func (m *mockStatusStore) GetNetwork(net *models.Network) error {
	return m.GetNetworkFunc(net)
}

func (m *mockStatusStore) GetSite(site *models.Site) error {
	return m.GetSiteFunc(site)
}

func TestGetStatus(t *testing.T) {
	app := echo.New()
	var s mockStore
	h := handler{Store: &s, Config: &models.Config{
		Timeout: 100,
	}}

	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))
	httpmock.RegisterResponder("GET", "https://err.test.site", httpmock.NewStringResponder(101, "fail"))
	httpmock.RegisterResponder("GET", "https://bigid.test.site", httpmock.NewStringResponder(200, "success"))
	httpmock.RegisterResponder("GET", "https://timeout.test.site", func(req *http.Request) (*http.Response, error) {
		time.Sleep(2 * time.Second)
		return httpmock.NewStringResponse(200, "success"), nil
	})

	ln, err := net.Listen("tcp", "[::]:12345")
	assert.NoError(t, err)
	defer ln.Close()

	cases := []struct {
		id      string
		uri     string
		isUp    bool
		wantErr error
	}{
		{id: "1", uri: "https://my.test.site", isUp: true, wantErr: nil},
		{id: "1", uri: "https://my.fail.site", isUp: false, wantErr: nil},
		{id: "1", uri: "https://^^invalidurl^^", isUp: false, wantErr: nil},
		{id: "1", uri: "ssh://localhost:12345", isUp: true, wantErr: nil},
		{id: "1", uri: "ssh://localhost:1234", isUp: false, wantErr: nil},
		{id: "1", uri: "https://500.test.site", isUp: false, wantErr: nil},
		{id: "abc", uri: "https://400.test.site", isUp: false, wantErr: echo.ErrBadRequest},
		{id: "", uri: "https://no-id.test.site", isUp: false, wantErr: echo.ErrBadRequest},
		{id: "12345", uri: "https://bigid.test.site", isUp: false, wantErr: echo.ErrNotFound},
		{id: "0", uri: "https://my.test.site", isUp: false, wantErr: echo.ErrBadRequest},
		{id: "1", uri: "https://timeout.test.site", isUp: false, wantErr: nil},
	}

	for _, c := range cases {
		s.GetSiteFunc = func(site *models.Site) error {
			if site.ID != 1 {
				return store.ErrNotFound
			}
			site.ID = 1
			site.URI = c.uri
			return nil
		}

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		ctx := app.NewContext(req, rec)
		ctx.SetPath("/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(c.id)
		err := h.getStatus(ctx)
		if c.wantErr != nil {
			assert.EqualError(t, err, c.wantErr.Error(), "%s should return %s", c.uri, c.wantErr.Error())
		} else {
			dec := json.NewDecoder(strings.NewReader(rec.Body.String()))
			ss := models.SiteStatus{}
			err := dec.Decode(&ss)
			assert.NoError(t, err)
			assert.Equal(t, c.isUp, ss.IsUp, "%s isUp should be %t", c.uri, c.isUp)
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
