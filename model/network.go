package model

import (
	"github.com/eiladin/go-simple-startpage/db"
)

type Network struct {
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
}

type Tag struct {
	Value string `json:"value"`
}

func LoadNetwork() Network {
	dbNetwork := db.ReadNetwork()
	return convertNetworkFromDbModel(dbNetwork)
}

func SaveNetwork(network Network) uint {
	dbNetwork := convertNetworkToDbModel(network)
	db.SaveNetwork(&dbNetwork)
	return dbNetwork.ID
}

func convertNetworkToDbModel(network Network) db.Network {
	dbNetwork := db.Network{
		Network: network.Network,
	}
	for _, link := range network.Links {
		dbLink := db.Link{
			Name:      link.Name,
			Uri:       link.Uri,
			SortOrder: link.SortOrder,
		}
		dbNetwork.Links = append(dbNetwork.Links, dbLink)
	}
	for _, site := range network.Sites {
		dbSite := db.Site{
			FriendlyName:   site.FriendlyName,
			Uri:            site.Uri,
			Icon:           site.Icon,
			IsSupportedApp: site.IsSupportedApp,
			SortOrder:      site.SortOrder,
		}
		for _, tag := range site.Tags {
			dbTag := db.Tag{
				Value: tag.Value,
			}
			dbSite.Tags = append(dbSite.Tags, dbTag)
		}
		dbNetwork.Sites = append(dbNetwork.Sites, dbSite)
	}
	return dbNetwork
}

func convertNetworkFromDbModel(dbNetwork db.Network) Network {
	network := Network{
		Network: dbNetwork.Network,
	}
	for _, dbLink := range dbNetwork.Links {
		link := Link{
			Name:      dbLink.Name,
			Uri:       dbLink.Uri,
			SortOrder: dbLink.SortOrder,
		}
		network.Links = append(network.Links, link)
	}
	for _, dbSite := range dbNetwork.Sites {
		site := Site{
			FriendlyName:   dbSite.FriendlyName,
			Uri:            dbSite.Uri,
			Icon:           dbSite.Icon,
			IsSupportedApp: dbSite.IsSupportedApp,
			SortOrder:      dbSite.SortOrder,
		}
		for _, dbTag := range dbSite.Tags {
			tag := Tag{
				Value: dbTag.Value,
			}
			site.Tags = append(site.Tags, tag)
		}
		network.Sites = append(network.Sites, site)
	}
	return network
}
