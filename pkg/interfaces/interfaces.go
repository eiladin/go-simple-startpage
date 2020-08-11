package interfaces

import (
	"time"

	"gorm.io/gorm"
)

// NetworkService interface
type NetworkService interface {
	CreateNetwork(net *Network)
	FindNetwork(net *Network)
	FindSite(site *Site)
}

// Network structure
type Network struct {
	ID        uint           `json:"-" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Network   string         `json:"network"`
	Links     []Link         `json:"links" gorm:"foreignkey:NetworkID"`
	Sites     []Site         `json:"sites" gorm:"foreignkey:NetworkID"`
}

// Link structure
type Link struct {
	ID        uint           `json:"-" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	NetworkID uint           `json:"-"`
	Name      string         `json:"name"`
	URI       string         `json:"uri"`
	SortOrder int            `json:"sortOrder"`
}

// Site structure
type Site struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	NetworkID      uint           `json:"-"`
	FriendlyName   string         `json:"friendlyName"`
	URI            string         `json:"uri"`
	Icon           string         `json:"icon"`
	IsSupportedApp bool           `json:"isSupportedApp"`
	SortOrder      int            `json:"sortOrder"`
	Tags           []Tag          `json:"tags" gorm:"foreignkey:SiteID"`
	IP             string         `json:"ip" gorm:"-"`
	IsUp           bool           `json:"isUp" gorm:"-"`
}

// Tag structure
type Tag struct {
	ID        uint           `json:"-" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	SiteID    uint           `json:"-"`
	Value     string         `json:"value"`
}
