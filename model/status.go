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
	config := LoadNetwork()
	return convertConfigToStatus(config)
}

func UpdateStatus(site StatusSite) *StatusSite {
	url, err := url.Parse(site.Uri)
	if err != nil {
		site.IsUp = false
		site.Ip = ""
		return &site
	}
	port := url.Port()
	if port == "22" {
		return testSSH(&site, url)
	} else {
		return testHTTP(&site, url)
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

func testSSH(site *StatusSite, url *url.URL) *StatusSite {
	address := url.Hostname() + ":" + "22"
	conn, err := net.Dial("tcp", address)
	if err != nil {
		site.IsUp = false
		site.Ip = ""
		return site
	}
	defer conn.Close()
	site.IsUp = true
	site.Ip = getIP(address)
	return site
}

func testHTTP(site *StatusSite, url *url.URL) *StatusSite {
	config := config.GetConfig()
	timeout := config.HealthCheck.Timeout
	sec := timeout / 1000
	dialer := &net.Dialer{
		Timeout: time.Duration(sec) * time.Second,
	}
	http.DefaultTransport.(*http.Transport).DialContext =
		func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, addr)
		}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, err := http.Get(site.Uri)
	if err != nil || resp.StatusCode < 200 || (resp.StatusCode >= 300 && resp.StatusCode != 401) {
		site.IsUp = false
		site.Ip = ""
		return site
	}
	defer resp.Body.Close()
	site.IsUp = true
	site.Ip = getIP(url.Host)

	return site
}

func convertConfigToStatus(network Network) Status {
	status := Status{
		Network: network.Network,
	}
	for _, link := range network.Links {
		statusLink := StatusLink{
			Name:      link.Name,
			Uri:       link.Uri,
			SortOrder: link.SortOrder,
		}
		status.Links = append(status.Links, statusLink)
	}
	for _, site := range network.Sites {
		statusSite := StatusSite{
			FriendlyName:   site.FriendlyName,
			Uri:            site.Uri,
			Icon:           site.Icon,
			IsSupportedApp: site.IsSupportedApp,
			SortOrder:      site.SortOrder,
			IsUp:           false,
			Ip:             "",
		}
		for _, tag := range site.Tags {
			statusTag := StatusTag{
				Value: tag.Value,
			}
			statusSite.Tags = append(statusSite.Tags, statusTag)
		}
		status.Sites = append(status.Sites, statusSite)
	}
	return status
}
