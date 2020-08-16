package model

import (
	"time"

	"gorm.io/gorm"
)

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
