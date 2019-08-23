package config

import (
	"sync"

	"github.com/eiladin/go-simple-startpage/db"
	"github.com/jinzhu/gorm"
)

var (
	once     sync.Once
	database *gorm.DB
)

const dbpath = "simple-startpage.db"

func init() {
	once.Do(initialize)
}

func initialize() {
	database = db.InitDB(dbpath)
	db.CreateDB(database)
}

type Config struct {
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

func Get() Config {
	dbConfig := db.ReadConfig(database)
	return convertFromDbModel(dbConfig)
}

func Add(config Config) uint {
	dbConfig := convertToDbModel(config)
	db.SaveConfig(database, &dbConfig)
	return dbConfig.ID
}

func convertToDbModel(config Config) db.Config {
	dbConfig := db.Config{}
	dbConfig.Network = config.Network
	for _, link := range config.Links {
		dbLink := db.Link{}
		dbLink.Name = link.Name
		dbLink.Uri = link.Uri
		dbLink.SortOrder = link.SortOrder
		dbConfig.Links = append(dbConfig.Links, dbLink)
	}
	for _, site := range config.Sites {
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
		dbConfig.Sites = append(dbConfig.Sites, dbSite)
	}
	return dbConfig
}

func convertFromDbModel(dbConfig db.Config) Config {
	config := Config{}
	config.Network = dbConfig.Network
	for _, dbLink := range dbConfig.Links {
		link := Link{}
		link.Name = dbLink.Name
		link.Uri = dbLink.Uri
		link.SortOrder = dbLink.SortOrder
		config.Links = append(config.Links, link)
	}
	for _, dbSite := range dbConfig.Sites {
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
		config.Sites = append(config.Sites, site)
	}
	return config
}
