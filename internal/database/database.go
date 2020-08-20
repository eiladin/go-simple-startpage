package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/internal/helpers"
	"github.com/eiladin/go-simple-startpage/pkg/model"
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
	conn, err := gorm.Open(dsn, cfg)
	if err != nil {
		helpers.Fatalf("Unable to connect to database: %v", err)
	}
	return conn
}

func getConfig(c *config.Config) *gorm.Config {
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

func getDSN(c *config.Config) gorm.Dialector {
	driver := strings.ToLower(c.Database.Driver)
	database := c.Database.Name
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

// MigrateDB runs database migrations
func MigrateDB(conn *gorm.DB) {
	conn.AutoMigrate(&model.Network{})
	conn.AutoMigrate(&model.Site{})
	conn.AutoMigrate(&model.Tag{})
	conn.AutoMigrate(&model.Link{})
}

// CreateNetwork creates a network in the database
func (d *DB) CreateNetwork(net *model.Network) {
	d.DB.Unscoped().Where("1 = 1").Delete(&model.Tag{})
	d.DB.Unscoped().Where("1 = 1").Delete(&model.Site{})
	d.DB.Unscoped().Where("1 = 1").Delete(&model.Link{})
	d.DB.Unscoped().Where("1 = 1").Delete(&model.Network{})
	d.DB.Create(&net)
}

// GetNetwork reads a network from the database
func (d *DB) GetNetwork(net *model.Network) {
	d.DB.Preload("Sites.Tags").Preload("Sites").Preload("Links").Find(net)
}

// GetSite reads a site from the database
func (d *DB) GetSite(site *model.Site) {
	d.DB.Find(site)
}
