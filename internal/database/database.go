package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiladin/go-simple-startpage/pkg/config"
	"github.com/eiladin/go-simple-startpage/pkg/network"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Compile-time proof of interface implementation.
var _ store.Store = (*DB)(nil)

var validDrivers = []string{"sqlite", "postgres", "mysql"}

type DB struct {
	conn *gorm.DB
}

type connectionRefusedErr string

func (e connectionRefusedErr) Error() string { return "unable to establish connection: " + string(e) }

type migrationFailedErr string

func (e migrationFailedErr) Error() string { return "unable to run database migrations: " + string(e) }

func New(config *config.Database) (store.Store, error) {
	d := DB{}
	cfg := getGormConfig(config)
	dsn := getDSN(config)
	conn, err := gorm.Open(dsn, cfg)
	if err != nil {
		return nil, connectionRefusedErr(err.Error())
	}
	d.conn = conn
	err = migrateDB(conn, config)
	if err != nil {
		return nil, migrationFailedErr(err.Error())
	}
	return &d, nil
}

func getGormConfig(c *config.Database) *gorm.Config {
	llevel := logger.Warn
	if c.Log {
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

func getDSN(c *config.Database) gorm.Dialector {
	driver := strings.ToLower(c.Driver)
	database := c.Name
	username := c.Username
	password := c.Password
	host := c.Host
	port := c.Port
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

func validDriver(driver string) bool {
	for _, d := range validDrivers {
		if d == driver {
			return true
		}
	}
	return false
}

func migrateDB(conn *gorm.DB, cfg *config.Database) error {
	if validDriver(cfg.Driver) {
		return conn.AutoMigrate(
			&network.Network{},
			&network.Site{},
			&network.DBTag{},
			&network.Link{},
		)
	} else {
		return fmt.Errorf("not a valid driver %s", cfg.Driver)
	}
}

func handleError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return store.ErrNotFound
	}
	return err
}

func (d *DB) CreateNetwork(net *network.Network) error {
	d.conn.Unscoped().Where("1 = 1").Delete(&network.DBTag{})
	d.conn.Unscoped().Where("1 = 1").Delete(&network.Site{})
	d.conn.Unscoped().Where("1 = 1").Delete(&network.Link{})
	d.conn.Unscoped().Where("1 = 1").Delete(&network.Network{})
	for i := range net.Sites {
		for _, tag := range net.Sites[i].Tags {
			net.Sites[i].DBTags = append(net.Sites[i].DBTags, network.DBTag{Value: tag})
		}
	}
	result := d.conn.Create(&net)
	return handleError(result.Error)
}

func (d *DB) GetNetwork(net *network.Network) error {
	result := d.conn.Preload("Sites.DBTags").Preload("Sites").Preload("Links").First(net)
	if result.Error != nil {
		return handleError(result.Error)
	}
	for i := range net.Sites {
		tags := []string{}
		for j := range net.Sites[i].DBTags {
			tags = append(tags, net.Sites[i].DBTags[j].Value)
		}
		net.Sites[i].Tags = tags
	}
	return nil
}

func (d *DB) GetSite(site *network.Site) error {
	result := d.conn.Where("name = ?", site.Name).First(site)
	return handleError(result.Error)
}

func (d *DB) Ping() error {
	sqlDB, _ := d.conn.DB()
	return sqlDB.Ping()
}
