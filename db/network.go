package db

import "github.com/jinzhu/gorm"

type Network struct {
	gorm.Model
	Network string
	Links   []Link `gorm:"foreignkey:NetworkID"`
	Sites   []Site `gorm:"foreignkey:NetworkID"`
}

type Link struct {
	gorm.Model
	NetworkID uint
	Name      string
	Uri       string
	SortOrder int
}

type Site struct {
	gorm.Model
	NetworkID      uint
	FriendlyName   string
	Uri            string
	Icon           string
	IsSupportedApp bool
	SortOrder      int
	Tags           []Tag `gorm:"foreignkey:SiteID"`
}

type Tag struct {
	gorm.Model
	SiteID uint
	Value  string
}

func MigrateDB() {
	db := GetDB()
	db.AutoMigrate(&Network{})
	db.AutoMigrate(&Site{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Link{})
}

func SaveNetwork(network *Network) {
	db := GetDB()
	db.Unscoped().Delete(Network{})
	db.Unscoped().Delete(Site{})
	db.Unscoped().Delete(Tag{})
	db.Unscoped().Delete(Link{})
	db.Create(&network)
}

func ReadNetwork() Network {
	db := GetDB()
	var result Network
	db.Set("gorm:auto_preload", true).Find(&result)
	return result
}

func DeleteNetwork(ID uint) {
	db := GetDB()
	db.Delete(Network{}, "ID = ?", ID)
}
