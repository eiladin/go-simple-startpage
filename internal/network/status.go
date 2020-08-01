package network

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/internal/database"
	"github.com/gofiber/fiber"
)

func updateStatus(s *Site) error {
	url, err := url.Parse(s.URI)
	if err != nil {
		return fmt.Errorf("Unable to parse URI: %s", s.URI)
	}
	port := url.Port()
	if port == "22" {
		return testSSH(s, url)
	}
	return testHTTP(s, url)
}

func getIP(host string) string {
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	if !strings.Contains(host, ".") {
		return host
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return ""
	}
	return ips[0].String()
}

func testSSH(s *Site, u *url.URL) error {
	address := u.Hostname() + ":" + "22"
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()
	s.IsUp = true
	s.IP = getIP(address)
	return nil
}

func testHTTP(s *Site, u *url.URL) error {
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
	s.IP = getIP(u.Host)
	return nil
}

// GetStatus handles /api/status
func GetStatus(c *fiber.Ctx) {
	db := database.DBConn
	var network Network
	db.Set("gorm:auto_preload", true).Find(&network)
	c.JSON(network)
}

// UpdateStatus handles /api/status
func UpdateStatus(c *fiber.Ctx) {
	var s Site
	err := c.BodyParser(&s)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return
	}
	err = updateStatus(&s)
	if err != nil {
		e := struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		}
		c.Status(fiber.StatusInternalServerError).JSON(e)
	}
	c.JSON(s)
}
