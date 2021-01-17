package models

import (
	"time"

	"gorm.io/gorm"
)

type Network struct {
	ID        uint           `json:"-" gorm:"primaryKey" yaml:"-"`
	CreatedAt time.Time      `json:"-" yaml:"-"`
	UpdatedAt time.Time      `json:"-" yaml:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" yaml:"-"`
	Network   string         `json:"network" yaml:"network"`
	Links     []Link         `json:"links" gorm:"foreignkey:NetworkID" yaml:"links"`
	Sites     []Site         `json:"sites" gorm:"foreignkey:NetworkID" yaml:"sites"`
}

type NetworkID struct {
	ID uint `json:"id" yaml:"id"`
}

type Link struct {
	ID        uint           `json:"-" gorm:"primaryKey" yaml:"-"`
	CreatedAt time.Time      `json:"-" yaml:"-"`
	UpdatedAt time.Time      `json:"-" yaml:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" yaml:"-"`
	NetworkID uint           `json:"-" yaml:"-"`
	Name      string         `json:"name" yaml:"name"`
	URI       string         `json:"uri" yaml:"uri"`
}

type Site struct {
	ID             uint           `json:"id" gorm:"primaryKey" yaml:"id"`
	CreatedAt      time.Time      `json:"-" yaml:"-"`
	UpdatedAt      time.Time      `json:"-" yaml:"-"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index" yaml:"-"`
	NetworkID      uint           `json:"-" yaml:"-"`
	FriendlyName   string         `json:"friendlyName" yaml:"friendlyName"`
	URI            string         `json:"uri" yaml:"uri"`
	Icon           string         `json:"icon" yaml:"icon"`
	IsSupportedApp bool           `json:"isSupportedApp" yaml:"isSupportedApp"`
	Tags           []Tag          `json:"tags" gorm:"foreignkey:SiteID" yaml:"tags"`
	IP             string         `json:"ip" gorm:"-" yaml:"-"`
	IsUp           bool           `json:"isUp" gorm:"-" yaml:"-"`
}

type Tag struct {
	ID        uint           `json:"-" gorm:"primaryKey" yaml:"-"`
	CreatedAt time.Time      `json:"-" yaml:"-"`
	UpdatedAt time.Time      `json:"-" yaml:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" yaml:"-"`
	SiteID    uint           `json:"-" yaml:"-"`
	Value     string         `json:"value" yaml:"value"`
}
