package handler

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/eiladin/go-simple-startpage/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

func updateStatus(timeout int, s *models.Site) error {
	url, err := url.Parse(s.URI)
	if err != nil {
		return fmt.Errorf("unable to parse URI: %s", s.URI)
	}
	s.IP = getIP(url)
	switch url.Scheme {
	case "ssh":
		err = testSSH(url)
	default:
		err = testHTTP(timeout, url)
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

func testHTTP(timeout int, u *url.URL) error {
	httpClient.Timeout = time.Millisecond * time.Duration(timeout)

	r, err := httpClient.Get(u.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode < 200 || (r.StatusCode >= 300 && r.StatusCode != 401) {
		return fmt.Errorf("invalid StatusCode: %d", r.StatusCode)
	}
	return nil
}

func getStatus(h handler, id uint) (*models.Site, error) {
	site := models.Site{ID: id}
	err := h.Store.GetSite(&site)
	if err != nil {
		return nil, err
	}
	err = updateStatus(h.Config.Timeout, &site)
	return &site, err
}

func (h handler) GetStatus(c echo.Context) error {
	val := c.Param("id")
	id, err := strconv.Atoi(val)
	if err != nil || id < 1 {
		if err == nil {
			err = errors.New("invalid id received: " + val)
		}
		return echo.ErrBadRequest.SetInternal(err)
	}
	site, err := getStatus(h, uint(id))
	if errors.Is(err, store.ErrNotFound) {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, models.SiteStatus{
		ID:   site.ID,
		IsUp: site.IsUp,
		IP:   site.IP,
	})
}

func (h handler) AddGetStatusRoute(app echoswagger.ApiRoot) echoswagger.ApiRoot {
	app.GET("/api/status/:id", h.GetStatus).
		AddParamPath(0, "id", "SiteID to get status for").
		AddResponse(http.StatusOK, "success", models.SiteStatus{}, nil).
		AddResponse(http.StatusBadRequest, "bad request", nil, nil).
		AddResponse(http.StatusNotFound, "not found", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)

	return app
}
