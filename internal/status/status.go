package status

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
)

// Handler handles Status commands
type Handler struct {
	SiteService interfaces.SiteService
}

func updateStatus(s *model.Site) error {
	url, err := url.Parse(s.URI)
	if err != nil {
		return fmt.Errorf("Unable to parse URI: %s", s.URI)
	}
	s.IP = getIP(url)
	if url.Scheme == "ssh" {
		return testSSH(s, url)
	}
	return testHTTP(s, url)
}

func getIP(u *url.URL) string {
	host := u.Hostname()
	ips, err := net.LookupIP(host)
	if err != nil {
		return ""
	}
	return ips[0].String()
}

func testSSH(s *model.Site, u *url.URL) error {
	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return err
	}
	defer conn.Close()
	s.IsUp = true
	return nil
}

var httpClient = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func testHTTP(s *model.Site, u *url.URL) error {
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
	s.IsUp = true
	return nil
}

func getStatus(h Handler, id uint) (*model.Site, error) {
	site := model.Site{ID: id}
	h.SiteService.FindSite(&site)
	err := updateStatus(&site)
	if err != nil {
		return &site, err
	}
	return &site, nil
}

// Get handles /api/status/:id
func (h Handler) Get(c echo.Context) error {
	val := c.Param("id")
	id, err := strconv.Atoi(val)
	if err != nil {
		return echo.ErrBadRequest
	}
	site, err := getStatus(h, uint(id))
	if err != nil {
		return echo.ErrInternalServerError
	}
	res := model.SiteStatus{
		ID:   site.ID,
		IsUp: site.IsUp,
		IP:   site.IP,
	}
	return c.JSON(http.StatusOK, res)
}
