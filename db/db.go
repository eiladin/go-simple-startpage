package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Config struct {
	gorm.Model
	Network string
	Links   []Link `gorm:"foreignkey:ConfigID"`
	Sites   []Site `gorm:"foreignkey:ConfigID"`
}

type Link struct {
	gorm.Model
	ConfigID  uint
	Name      string
	Uri       string
	SortOrder int
}

type Site struct {
	gorm.Model
	ConfigID       uint
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

func InitDB(filepath string) *gorm.DB {
	db, err := gorm.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	return db
}

func CreateDB(db *gorm.DB) {
	db.AutoMigrate(&Config{})
	db.AutoMigrate(&Site{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Link{})
}

func SaveConfig(db *gorm.DB, config *Config) {
	db.Unscoped().Delete(Config{})
	db.Unscoped().Delete(Site{})
	db.Unscoped().Delete(Tag{})
	db.Unscoped().Delete(Link{})
	db.Create(&config)
}

func ReadConfig(db *gorm.DB) Config {
	var result Config
	db.Set("gorm:auto_preload", true).Find(&result)
	return result
}

func DeleteConfig(db *gorm.DB, ID uint) {
	db.Delete(Config{}, "ID = ?", ID)
}
