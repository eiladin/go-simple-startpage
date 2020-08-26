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

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/internal/store"
	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/labstack/echo/v4"
	"github.com/pangpanglabs/echoswagger/v2"
)

// Status struct
type Status struct {
	Store store.Store
}

func updateStatus(s *model.Site) error {
	url, err := url.Parse(s.URI)
	if err != nil {
		return fmt.Errorf("unable to parse URI: %s", s.URI)
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
		return fmt.Errorf("invalid StatusCode: %d", r.StatusCode)
	}
	return nil
}

func getStatus(h Status, id uint) (*model.Site, error) {
	site := model.Site{ID: id}
	err := h.Store.GetSite(&site)
	if err != nil {
		return nil, err
	}
	err = updateStatus(&site)
	return &site, err
}

// Get /api/status/{id}
func (h Status) Get(c echo.Context) error {
	val := c.Param("id")
	id, err := strconv.Atoi(val)
	if err != nil || id < 1 {
		if err == nil {
			err = errors.New("invalid id received: " + val)
		}
		return echo.ErrBadRequest.SetInternal(err)
	}
	site, err := getStatus(h, uint(id))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return echo.ErrNotFound
		}
		return echo.ErrInternalServerError.SetInternal(err)
	}
	res := model.SiteStatus{
		ID:   site.ID,
		IsUp: site.IsUp,
		IP:   site.IP,
	}
	return c.JSON(http.StatusOK, res)
}

// Register handler
func (h Status) Register(app echoswagger.ApiRoot) echoswagger.ApiRoot {
	app.GET("/api/status/:id", h.Get).
		AddParamPath(0, "id", "SiteID to get status for").
		AddResponse(http.StatusOK, "success", model.SiteStatus{}, nil).
		AddResponse(http.StatusBadRequest, "bad request", nil, nil).
		AddResponse(http.StatusNotFound, "not found", nil, nil).
		AddResponse(http.StatusInternalServerError, "internal server error", nil, nil)

	return app
}
