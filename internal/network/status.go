package network

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/gofiber/fiber"
)

func (s *Site) updateStatus() error {
	url, err := url.Parse(s.URI)
	if err != nil {
		return fmt.Errorf("Unable to parse URI: %s", s.URI)
	}
	s.IP = getIP(url)
	if url.Scheme == "ssh" {
		return s.testSSH(url)
	}
	return s.testHTTP(url)
}

func getIP(u *url.URL) string {
	host := u.Hostname()
	ips, err := net.LookupIP(host)
	if err != nil {
		return ""
	}
	return ips[0].String()
}

func (s *Site) testSSH(u *url.URL) error {
	conn, err := net.Dial("tcp", u.Host)
	if err != nil {
		return err
	}
	defer conn.Close()
	s.IsUp = true
	return nil
}

func (s *Site) testHTTP(u *url.URL) error {
	c := config.GetConfig()
	dialer := &net.Dialer{
		Timeout: time.Duration(c.HealthCheck.Timeout) * time.Millisecond,
	}
	http.DefaultTransport.(*http.Transport).DialContext =
		func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	r, err := http.Get(s.URI)
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
func UpdateStatus(c *fiber.Ctx) {
	var s Site
	err := c.BodyParser(&s)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return
	}
	err = s.updateStatus()
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.JSON(s)
}
