package models

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Status struct {
	Name string `json:"name"`
	IsUp bool   `json:"isUp"`
	IP   string `json:"ip"`
}

func NewStatus(timeout int, s *Site) Status {
	res := Status{
		Name: s.Name,
		IsUp: false,
	}
	url, err := url.Parse(s.URI)
	if err != nil {
		return res
	}

	res.IP = getIP(url)
	switch url.Scheme {
	case "ssh":
		err = testSSH(url)
	default:
		err = testHTTP(timeout, url)
	}
	res.IsUp = err == nil
	return res
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
