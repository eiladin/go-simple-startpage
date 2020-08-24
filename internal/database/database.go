package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/eiladin/go-simple-startpage/internal/store"
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

type connectionRefusedErr string

func (e connectionRefusedErr) Error() string { return "Unable to establish connection: " + string(e) }

// InitDB initialized the selected database
func InitDB() (*gorm.DB, error) {
	c := config.GetConfig()
	cfg := getConfig(&c)
	dsn := getDSN(&c)
	conn, err := gorm.Open(dsn, cfg)
	if err != nil {
		return nil, connectionRefusedErr(err.Error())
	}
	return conn, nil
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

func handleError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return store.ErrNotFound
	}
	return err
}

// CreateNetwork creates a network in the database
func (d *DB) CreateNetwork(net *model.Network) error {
	d.DB.Unscoped().Where("1 = 1").Delete(&model.Tag{})
	d.DB.Unscoped().Where("1 = 1").Delete(&model.Site{})
	d.DB.Unscoped().Where("1 = 1").Delete(&model.Link{})
	d.DB.Unscoped().Where("1 = 1").Delete(&model.Network{})
	result := d.DB.Create(&net)
	return handleError(result.Error)
}

// GetNetwork reads a network from the database
func (d *DB) GetNetwork(net *model.Network) error {
	result := d.DB.Preload("Sites.Tags").Preload("Sites").Preload("Links").First(net)
	return handleError(result.Error)
}

// GetSite reads a site from the database
func (d *DB) GetSite(site *model.Site) error {
	result := d.DB.First(site)
	return handleError(result.Error)
}
