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
