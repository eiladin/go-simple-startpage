package model

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
)

type SiteStatus struct {
	ID   uint   `json:"id"`
	IsUp bool   `json:"isUp"`
	IP   string `json:"ip"`
}

func NewSiteStatus(httpClient http.Client, s *Site) SiteStatus {
	res := SiteStatus{
		ID:   s.ID,
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
		err = testHTTP(httpClient, url)
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

func testHTTP(httpClient http.Client, u *url.URL) error {
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
