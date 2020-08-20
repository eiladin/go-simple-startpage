package handler

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
)

// Status struct
type Status struct {
	Store Store
}

func updateStatus(s *model.Site) error {
	url, err := url.Parse(s.URI)
	if err != nil {
		return fmt.Errorf("Unable to parse URI: %s", s.URI)
	}
	s.IP = getIP(url)
	switch url.Scheme {
	case "ssh":
		err = testSSH(url)
	default:
		err = testHTTP(url)
	}
	s.IsUp = err == nil
	return err
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

func testHTTP(u *url.URL) error {
	c := config.GetConfig()
	httpClient.Timeout = time.Millisecond * time.Duration(c.Timeout)

	r, err := httpClient.Get(u.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode < 200 || (r.StatusCode >= 300 && r.StatusCode != 401) {
		return fmt.Errorf("Invalid StatusCode: %d", r.StatusCode)
	}
	return nil
}

func getStatus(h Status, id uint) (*model.Site, error) {
	site := model.Site{ID: id}
	h.Store.GetSite(&site)
	err := updateStatus(&site)
	return &site, err
}

// Get handles /api/status/:id
func (h Status) Get(c echo.Context) error {
	val := c.Param("id")
	id, err := strconv.Atoi(val)
	if err != nil {
		return echo.ErrBadRequest
	}
	site, _ := getStatus(h, uint(id))
	res := model.SiteStatus{
		ID:   site.ID,
		IsUp: site.IsUp,
		IP:   site.IP,
	}
	return c.JSON(http.StatusOK, res)
}
