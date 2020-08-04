package network

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
	"github.com/gofiber/fiber"
)

func updateStatus(s *interfaces.Site) error {
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

func testSSH(s *interfaces.Site, u *url.URL) error {
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

func testHTTP(s *interfaces.Site, u *url.URL) error {
	c := config.GetConfig()
	httpClient.Timeout = time.Millisecond * time.Duration(c.HealthCheck.Timeout)

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

// UpdateStatus handles /api/status
func (h Handler) UpdateStatus(c *fiber.Ctx) {
	var s interfaces.Site
	err := c.BodyParser(&s)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return
	}
	err = updateStatus(&s)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.JSON(s)
}
