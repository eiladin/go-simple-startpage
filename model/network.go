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
	database := db.GetDB()
	dbNetwork := db.ReadNetwork(database)
	return convertNetworkFromDbModel(dbNetwork)
}

func SaveNetwork(network Network) uint {
	database := db.GetDB()
	dbNetwork := convertNetworkToDbModel(network)
	db.SaveNetwork(database, &dbNetwork)
	return dbNetwork.ID
}

func convertNetworkToDbModel(network Network) db.Network {
	dbNetwork := db.Network{}
	dbNetwork.Network = network.Network
	for _, link := range network.Links {
		dbLink := db.Link{}
		dbLink.Name = link.Name
		dbLink.Uri = link.Uri
		dbLink.SortOrder = link.SortOrder
		dbNetwork.Links = append(dbNetwork.Links, dbLink)
	}
	for _, site := range network.Sites {
		dbSite := db.Site{}
		dbSite.FriendlyName = site.FriendlyName
		dbSite.Uri = site.Uri
		dbSite.Icon = site.Icon
		dbSite.IsSupportedApp = site.IsSupportedApp
		dbSite.SortOrder = site.SortOrder
		for _, tag := range site.Tags {
			dbTag := db.Tag{}
			dbTag.Value = tag.Value
			dbSite.Tags = append(dbSite.Tags, dbTag)
		}
		dbNetwork.Sites = append(dbNetwork.Sites, dbSite)
	}
	return dbNetwork
}

func convertNetworkFromDbModel(dbNetwork db.Network) Network {
	network := Network{}
	network.Network = dbNetwork.Network
	for _, dbLink := range dbNetwork.Links {
		link := Link{}
		link.Name = dbLink.Name
		link.Uri = dbLink.Uri
		link.SortOrder = dbLink.SortOrder
		network.Links = append(network.Links, link)
	}
	for _, dbSite := range dbNetwork.Sites {
		site := Site{}
		site.FriendlyName = dbSite.FriendlyName
		site.Uri = dbSite.Uri
		site.Icon = dbSite.Icon
		site.IsSupportedApp = dbSite.IsSupportedApp
		site.SortOrder = dbSite.SortOrder
		for _, dbTag := range dbSite.Tags {
			tag := Tag{}
			tag.Value = dbTag.Value
			site.Tags = append(site.Tags, tag)
		}
		network.Sites = append(network.Sites, site)
	}
	return network
}
