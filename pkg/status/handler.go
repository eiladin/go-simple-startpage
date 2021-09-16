package status

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/network"
)

var (
	ErrNotFound = errors.New("Status not found")
)

type IHandler interface {
	Get(string) (*Status, error)
}

type repository interface {
	GetSite(*network.Site) error
}

// Compile-time proof of interface implementation.
var _ IHandler = (*handler)(nil)

type handler struct {
	repo   repository
	config *config.Config
}

func NewHandler(repo repository, cfg *config.Config) IHandler {
	return &handler{
		repo:   repo,
		config: cfg,
	}
}

func (c *handler) Get(name string) (*Status, error) {
	site := network.Site{Name: name}

	if err := c.repo.GetSite(&site); err != nil {
		return nil, ErrNotFound
	}

	return checkSiteStatus(c.config.Timeout, site), nil
}

func checkSiteStatus(timeout int, site network.Site) *Status {
	url, err := url.Parse(site.URI)
	if err != nil {
		return &Status{Name: site.Name}
	}

	return &Status{
		Name: site.Name,
		IsUp: testConnection(timeout, url),
		IP:   getIP(url),
	}
}

func testConnection(timeout int, url *url.URL) bool {
	if url.Scheme == "ssh" {
		return testSSH(url) == nil
	} else {
		return testHTTP(timeout, url) == nil
	}
}

func getIP(u *url.URL) string {
	host := u.Hostname()
	ips, err := net.LookupIP(host)
	if err != nil {
		return ""
	}
	return ips[0].String()
}

func testSSH(u *url.URL) error {
	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

var httpClient = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func testHTTP(timeout int, u *url.URL) error {
	httpClient.Timeout = time.Millisecond * time.Duration(timeout)
	r, err := httpClient.Get(u.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode < 200 || (r.StatusCode >= 300 && r.StatusCode != 401) {
		return fmt.Errorf("invalid status: %d", r.StatusCode)
	}
	return nil
}
