package models

import (
	"time"

	"gorm.io/gorm"
)

type Network struct {
	ID        uint           `json:"-" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Network   string         `json:"network"`
	Links     []Link         `json:"links" gorm:"foreignkey:NetworkID"`
	Sites     []Site         `json:"sites" gorm:"foreignkey:NetworkID"`
}

type NetworkID struct {
	ID uint `json:"id"`
}

type Link struct {
	ID        uint           `json:"-" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	NetworkID uint           `json:"-"`
	Name      string         `json:"name"`
	URI       string         `json:"uri"`
}

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
	Tags           []Tag          `json:"tags" gorm:"foreignkey:SiteID"`
	IP             string         `json:"ip" gorm:"-"`
	IsUp           bool           `json:"isUp" gorm:"-"`
}

type Tag struct {
	ID        uint           `json:"-" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	SiteID    uint           `json:"-"`
	Value     string         `json:"value"`
}
