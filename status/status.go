package status

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
	Network string `json:"network"`
	Links   []Link `json:"links"`
	Sites   []Site `json:"sites"`
}

type Link struct {
	Name      string `json:"name"`
	Uri       string `json:"uri"`
	SortOrder int    `json:"sortOrder"`
}

type Site struct {
	FriendlyName   string `json:"friendlyName"`
	Uri            string `json:"uri"`
	Icon           string `json:"icon"`
	IsSupportedApp bool   `json:"isSupportedApp"`
	SortOrder      int    `json:"sortOrder"`
	Tags           []Tag  `json:"tags"`
	Ip             string `json:"ip"`
	IsUp           bool   `json:"isUp"`
}

type Tag struct {
	Value string `json:"value"`
}

func Get() Status {
	config := config.Get()
	return convertConfigToStatus(config)
}

func UpdateStatus(site Site) *Site {
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

func testSSH(site *Site, url *url.URL) *Site {
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

func testHTTP(site *Site, url *url.URL) *Site {
	dialer := &net.Dialer{
		Timeout: 2 * time.Second,
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

func convertConfigToStatus(config config.Config) Status {
	status := Status{}
	status.Network = config.Network
	for _, link := range config.Links {
		statusLink := Link{}
		statusLink.Name = link.Name
		statusLink.Uri = link.Uri
		statusLink.SortOrder = link.SortOrder
		status.Links = append(status.Links, statusLink)
	}
	for _, site := range config.Sites {
		statusSite := Site{}
		statusSite.FriendlyName = site.FriendlyName
		statusSite.Uri = site.Uri
		statusSite.Icon = site.Icon
		statusSite.IsSupportedApp = site.IsSupportedApp
		statusSite.SortOrder = site.SortOrder
		statusSite.IsUp = false
		statusSite.Ip = ""
		for _, tag := range site.Tags {
			statusTag := Tag{}
			statusTag.Value = tag.Value
			statusSite.Tags = append(statusSite.Tags, statusTag)
		}
		status.Sites = append(status.Sites, statusSite)
	}
	return status
}
