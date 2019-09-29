package model

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/eiladin/go-simple-startpage/config"
)

type Status struct {
	Network string       `json:"network"`
	Links   []StatusLink `json:"links"`
	Sites   []StatusSite `json:"sites"`
}

type StatusLink struct {
	Name      string `json:"name"`
	Uri       string `json:"uri"`
	SortOrder int    `json:"sortOrder"`
}

type StatusSite struct {
	FriendlyName   string      `json:"friendlyName"`
	Uri            string      `json:"uri"`
	Icon           string      `json:"icon"`
	IsSupportedApp bool        `json:"isSupportedApp"`
	SortOrder      int         `json:"sortOrder"`
	Tags           []StatusTag `json:"tags"`
	Ip             string      `json:"ip"`
	IsUp           bool        `json:"isUp"`
}

type StatusTag struct {
	Value string `json:"value"`
}

func GetStatus() Status {
	n := LoadNetwork()
	return convertNetworkToStatus(n)
}

func UpdateStatus(s StatusSite) *StatusSite {
	url, err := url.Parse(s.Uri)
	if err != nil {
		s.IsUp = false
		s.Ip = ""
		return &s
	}
	port := url.Port()
	if port == "22" {
		return testSSH(&s, url)
	} else {
		return testHTTP(&s, url)
	}
}

func getIP(host string) string {
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	if !strings.Contains(host, ".") {
		return host
	} else {
		ips, err := net.LookupIP(host)
		if err != nil {
			return ""
		} else {
			return ips[0].String()
		}
	}
}

func testSSH(s *StatusSite, u *url.URL) *StatusSite {
	address := u.Hostname() + ":" + "22"
	conn, err := net.Dial("tcp", address)
	if err != nil {
		s.IsUp = false
		s.Ip = ""
		return s
	}
	defer conn.Close()
	s.IsUp = true
	s.Ip = getIP(address)
	return s
}

func testHTTP(s *StatusSite, u *url.URL) *StatusSite {
	c := config.GetConfig()
	timeout := c.HealthCheck.Timeout
	sec := timeout / 1000
	dialer := &net.Dialer{
		Timeout: time.Duration(sec) * time.Second,
	}
	http.DefaultTransport.(*http.Transport).DialContext =
		func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	r, err := http.Get(s.Uri)
	if err != nil || r.StatusCode < 200 || (r.StatusCode >= 300 && r.StatusCode != 401) {
		s.IsUp = false
		s.Ip = ""
		return s
	}
	defer r.Body.Close()
	s.IsUp = true
	s.Ip = getIP(u.Host)

	return s
}

func convertNetworkToStatus(n Network) Status {
	s := Status{
		Network: n.Network,
	}
	for _, l := range n.Links {
		sl := StatusLink{
			Name:      l.Name,
			Uri:       l.Uri,
			SortOrder: l.SortOrder,
		}
		s.Links = append(s.Links, sl)
	}
	for _, st := range n.Sites {
		ss := StatusSite{
			FriendlyName:   st.FriendlyName,
			Uri:            st.Uri,
			Icon:           st.Icon,
			IsSupportedApp: st.IsSupportedApp,
			SortOrder:      st.SortOrder,
			IsUp:           false,
			Ip:             "",
		}
		for _, tag := range st.Tags {
			st := StatusTag{
				Value: tag.Value,
			}
			ss.Tags = append(ss.Tags, st)
		}
		s.Sites = append(s.Sites, ss)
	}
	return s
}
