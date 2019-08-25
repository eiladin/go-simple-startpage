package db

import (
	"fmt"

	"github.com/eiladin/go-simple-startpage/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

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

var DB *gorm.DB

func InitDB() *gorm.DB {
	var db *gorm.DB
	var err error
	config := config.InitConfig()

	driver := config.Database.Driver
	database := config.Database.Dbname
	username := config.Database.Username
	password := config.Database.Password
	host := config.Database.Host
	port := config.Database.Port

	if driver == "sqlite" {
		db, err = gorm.Open("sqlite3", database)
		if err != nil {
			fmt.Println("db err: ", err)
		}
	} else if driver == "postgres" {
		db, err = gorm.Open("postgres", "host="+host+" port="+port+" user="+username+" dbname="+database+"  sslmode=disable password="+password)
		if err != nil {
			fmt.Println("db err: ", err)
		}
	} else if driver == "mysql" {
		db, err = gorm.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/"+database+"?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			fmt.Println("db err: ", err)
		}
	} else {
		panic("database driver is undefined")
	}

	db.LogMode(true)
	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}

func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(&Network{})
	db.AutoMigrate(&Site{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Link{})
}

func SaveNetwork(db *gorm.DB, network *Network) {
	db.Unscoped().Delete(Network{})
	db.Unscoped().Delete(Site{})
	db.Unscoped().Delete(Tag{})
	db.Unscoped().Delete(Link{})
	db.Create(&network)
}

func ReadNetwork(db *gorm.DB) Network {
	var result Network
	db.Set("gorm:auto_preload", true).Find(&result)
	return result
}

func DeleteNetwork(db *gorm.DB, ID uint) {
	db.Delete(Network{}, "ID = ?", ID)
}
