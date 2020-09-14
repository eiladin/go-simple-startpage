package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiladin/go-simple-startpage/pkg/model"
	"github.com/eiladin/go-simple-startpage/pkg/store"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	conn *gorm.DB
}

type connectionRefusedErr string

func (e connectionRefusedErr) Error() string { return "unable to establish connection: " + string(e) }

type migrationFailedErr string

func (e migrationFailedErr) Error() string { return "unable to run database migrations: " + string(e) }

func New(config *model.Config) (store.Store, error) {
	d := DB{}
	cfg := getGormConfig(config)
	dsn := getDSN(config)
	conn, err := gorm.Open(dsn, cfg)
	if err != nil {
		return nil, connectionRefusedErr(err.Error())
	}
	d.conn = conn
	err = migrateDB(conn)
	if err != nil {
		return nil, migrationFailedErr(err.Error())
	}
	return &d, nil
}

func getGormConfig(c *model.Config) *gorm.Config {
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

func getDSN(c *model.Config) gorm.Dialector {
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

func migrateDB(conn *gorm.DB) error {
	return conn.AutoMigrate(
		&model.Network{},
		&model.Site{},
		&model.Tag{},
		&model.Link{},
	)
}

func handleError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return store.ErrNotFound
	}
	return err
}

func (d *DB) CreateNetwork(net *model.Network) error {
	d.conn.Unscoped().Where("1 = 1").Delete(&model.Tag{})
	d.conn.Unscoped().Where("1 = 1").Delete(&model.Site{})
	d.conn.Unscoped().Where("1 = 1").Delete(&model.Link{})
	d.conn.Unscoped().Where("1 = 1").Delete(&model.Network{})
	result := d.conn.Create(&net)
	return handleError(result.Error)
}

func (d *DB) GetNetwork(net *model.Network) error {
	result := d.conn.Preload("Sites.Tags").Preload("Sites").Preload("Links").First(net)
	return handleError(result.Error)
}

func (d *DB) GetSite(site *model.Site) error {
	result := d.conn.First(site)
	return handleError(result.Error)
}

func (d *DB) Ping() error {
	sqlDB, err := d.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
