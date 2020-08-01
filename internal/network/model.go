package network

import "gorm.io/gorm"

// Network structure
type Network struct {
	gorm.Model
	Network string `json:"network"`
	Links   []Link `json:"links" gorm:"foreignkey:NetworkID"`
	Sites   []Site `json:"sites" gorm:"foreignkey:NetworkID"`
}

// Link structure
type Link struct {
	gorm.Model
	NetworkID uint
	Name      string `json:"name"`
	URI       string `json:"uri"`
	SortOrder int    `json:"sortOrder"`
}

// Site structure
type Site struct {
	gorm.Model
	NetworkID      uint
	FriendlyName   string `json:"friendlyName"`
	URI            string `json:"uri"`
	Icon           string `json:"icon"`
	IsSupportedApp bool   `json:"isSupportedApp"`
	SortOrder      int    `json:"sortOrder"`
	Tags           []Tag  `json:"tags" gorm:"foreignkey:SiteID"`
	IP             string `json:"ip" gorm:"-"`
	IsUp           bool   `json:"isUp" gorm:"-"`
}

// Tag structure
type Tag struct {
	gorm.Model
	SiteID uint
	Value  string `json:"value"`
}
