package status

import "github.com/eiladin/go-simple-startpage/config"

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

func UpdateStatus(site Site) Site {
	site.IsUp = true
	site.Ip = "0.0.0.0"
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
