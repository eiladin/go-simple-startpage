package model

import (
	"time"

	"gorm.io/gorm"
)

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
