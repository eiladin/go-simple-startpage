package models

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type SiteStatusSuite struct {
	suite.Suite
}

func (suite *SiteStatusSuite) TestHttp() {
	cases := []struct {
		url     string
		timeout int
		wantErr bool
	}{
		{url: "https://my.test.site", timeout: 0, wantErr: false},
		{url: "https://timeout.test.site", timeout: 100, wantErr: true},
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))
	httpmock.RegisterResponder("GET", "https://timeout.test.site", func(req *http.Request) (*http.Response, error) {
		time.Sleep(2 * time.Second)
		return httpmock.NewStringResponse(200, "success"), nil
	})

	for _, c := range cases {
		httpClient.Timeout = time.Millisecond * time.Duration(c.timeout)
		url, err := url.Parse(c.url)
		suite.NoError(err)
		err = testHTTP(httpClient, url)
		if c.wantErr {
			suite.Error(err)
		} else {
			suite.NoError(err)
		}
	}
}

func (suite *SiteStatusSuite) TestTCP() {
	ln, err := net.Listen("tcp", "[::]:22222")
	suite.NoError(err)
	defer ln.Close()

	url, err := url.Parse("ssh://localhost:22222")
	suite.NoError(err)
	err = testSSH(url)
	suite.NoError(err, "ssh://localhost:22222 should not error")
}

func (suite *SiteStatusSuite) TestGetIP() {
	url, err := url.Parse("http://localhost")
	suite.NoError(err)
	ip := getIP(url)
	suite.Contains([]string{"127.0.0.1", "::1"}, ip, "http://localhost should return the following ips: [127.0.0.1, ::1]")
}

func (suite *SiteStatusSuite) TestNewSiteStatus() {
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

	httpClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	httpmock.ActivateNonDefault(&httpClient)
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://my.test.site", httpmock.NewStringResponder(200, "success"))
	httpmock.RegisterResponder("GET", "https://err.test.site", httpmock.NewStringResponder(101, "redirect"))

	ln, err := net.Listen("tcp", "[::]:22223")
	suite.NoError(err)
	defer ln.Close()

	for _, c := range cases {
		s := NewSiteStatus(httpClient, &c.site)
		suite.Equal(c.isUp, s.IsUp, "site %s isUp should be %t", c.site.URI, c.isUp)
	}
}

func TestSiteStatusSuite(t *testing.T) {
	suite.Run(t, new(SiteStatusSuite))
}
