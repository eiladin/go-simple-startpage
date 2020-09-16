package models

import (
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type StatusSuite struct {
	suite.Suite
}

func (suite *StatusSuite) TestHttp() {
	cases := []struct {
		url     string
		timeout int
		wantErr bool
	}{
		{url: "https://my.test.site", timeout: 0, wantErr: false},
		{url: "https://timeout.test.site", timeout: 100, wantErr: true},
	}

	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))
	httpmock.RegisterResponder("GET", "https://timeout.test.site", func(req *http.Request) (*http.Response, error) {
		time.Sleep(2 * time.Second)
		return httpmock.NewStringResponse(200, "success"), nil
	})

	for _, c := range cases {
		url, err := url.Parse(c.url)
		suite.NoError(err)
		err = testHTTP(c.timeout, url)
		if c.wantErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
	}
}

func (suite *StatusSuite) TestTCP() {
	ln, err := net.Listen("tcp", "[::]:22222")
	suite.NoError(err)
	defer ln.Close()

	url, err := url.Parse("ssh://localhost:22222")
	suite.NoError(err)
	err = testSSH(url)
	suite.NoError(err, "ssh://localhost:22222 should not error")
}

func (suite *StatusSuite) TestGetIP() {
	url, err := url.Parse("http://localhost")
	suite.NoError(err)
	ip := getIP(url)
	suite.Contains([]string{"127.0.0.1", "::1"}, ip, "http://localhost should return the following ips: [127.0.0.1, ::1]")
}

func (suite *StatusSuite) TestNewStatus() {
	cases := []struct {
		site Site
		isUp bool
	}{
		{site: Site{URI: "https://my.test.site"}, isUp: true},
		{site: Site{URI: "https://my.fail.site"}, isUp: false},
		{site: Site{URI: "https://^^invalidurl^^"}, isUp: false},
		{site: Site{URI: "ssh://localhost:22223"}, isUp: true},
		{site: Site{URI: "ssh://localhost:1234"}, isUp: false},
		{site: Site{URI: "https://err.test.site"}, isUp: false},
	}

	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))
	httpmock.RegisterResponder("GET", "https://err.test.site", httpmock.NewStringResponder(101, "redirect"))

	ln, err := net.Listen("tcp", "[::]:22223")
	suite.NoError(err)
	defer ln.Close()

	for _, c := range cases {
		s := NewStatus(0, &c.site)
		suite.Equal(c.isUp, s.IsUp, "site %s isUp should be %t", c.site.URI, c.isUp)
	}
}

func TestStatusSuite(t *testing.T) {
	suite.Run(t, new(StatusSuite))
}
