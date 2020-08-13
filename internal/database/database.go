package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/pkg/interfaces"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB structure
type DB struct {
	DB *gorm.DB
}

// InitDB initialized the selected database
func InitDB() *gorm.DB {
	c := config.GetConfig()
	cfg := getConfig(&c)
	dsn := getDSN(&c)
	return open(dsn, cfg)
}

func getConfig(c *config.Configuration) *gorm.Config {
	llevel := logger.Silent
	if c.Database.Log {
		llevel = logger.Info
	}
	return &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: llevel,
			},
		),
	}
}

func getDSN(c *config.Configuration) gorm.Dialector {
	driver := strings.ToLower(c.Database.Driver)
	database := c.Database.Dbname
	username := c.Database.Username
	password := c.Database.Password
	host := c.Database.Host
	port := c.Database.Port
	if driver == "sqlite" {
		return sqlite.Open(database)
	} else if driver == "postgres" {
		dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, username, database, password)
		return postgres.Open(dsn)
	} else if driver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)
		return mysql.Open(dsn)
	}
	return sqlite.Open("simple-startpage.db")

}

func open(dialector gorm.Dialector, config *gorm.Config) *gorm.DB {
	conn, err := gorm.Open(dialector, config)
	if err != nil {
		fmt.Println("db err: ", err)
	}
	return conn
}

// CreateNetwork creates a network in the database
func (d *DB) CreateNetwork(net *interfaces.Network) {
	d.DB.Unscoped().Where("1 = 1").Delete(&interfaces.Tag{})
	d.DB.Unscoped().Where("1 = 1").Delete(&interfaces.Site{})
	d.DB.Unscoped().Where("1 = 1").Delete(&interfaces.Link{})
	d.DB.Unscoped().Where("1 = 1").Delete(&interfaces.Network{})
	d.DB.Create(&net)
}

// FindNetwork reads a network from the database
func (d *DB) FindNetwork(net *interfaces.Network) {
	d.DB.Preload("Sites.Tags").Preload("Sites").Preload("Links").Find(net)
}

// FindSite reads a site from the database
func (d *DB) FindSite(site *interfaces.Site) {
	d.DB.Find(site)
}
