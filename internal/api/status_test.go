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
	"github.com/stretchr/testify/suite"
)

type StatusServiceSuite struct {
	suite.Suite
}

func (suite *StatusServiceSuite) TestGet() {
	app := echo.New()

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
	suite.NoError(err)
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

		cfg := &models.Config{Timeout: 100}
		s := &mockStore{
			GetSiteFunc: func(site *models.Site) error {
				if site.ID != 1 {
					return store.ErrNotFound
				}
				site.ID = 1
				site.URI = c.uri
				return nil
			},
		}
		ss := NewStatusService(cfg, s)

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		ctx := app.NewContext(req, rec)
		ctx.SetPath("/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(c.id)
		err := ss.Get(ctx)
		if c.wantErr != nil {
			suite.EqualError(err, c.wantErr.Error(), "%s should return %s", c.uri, c.wantErr.Error())
		} else {
			dec := json.NewDecoder(strings.NewReader(rec.Body.String()))
			ss := models.SiteStatus{}
			err := dec.Decode(&ss)
			suite.NoError(err)
			suite.Equal(c.isUp, ss.IsUp, "%s isUp should be %t", c.uri, c.isUp)
		}
	}
}

func (suite *StatusServiceSuite) TestRegister() {
	app := echoswagger.New(echo.New(), "/swagger-test", &echoswagger.Info{})
	NewStatusService(&models.Config{}, &mockStore{}).Register(app)
	e := []string{}
	for _, r := range app.Echo().Routes() {
		e = append(e, r.Method+" "+r.Path)
	}
	suite.Contains(e, "GET /api/status/:id")
}

func TestStatusServiceSuite(t *testing.T) {
	suite.Run(t, new(StatusServiceSuite))
}
