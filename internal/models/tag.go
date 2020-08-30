package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	ID        uint           `json:"-" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	SiteID    uint           `json:"-"`
	Value     string         `json:"value"`
}
