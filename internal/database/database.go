package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBConn holds the database connection
var DBConn *gorm.DB

// InitDB initialized the selected database
func InitDB() {
	var err error
	c := config.GetConfig()

	driver := strings.ToLower(c.Database.Driver)
	database := c.Database.Dbname
	username := c.Database.Username
	password := c.Database.Password
	host := c.Database.Host
	port := c.Database.Port
	llevel := logger.Silent
	if c.Database.Log {
		llevel = logger.Info
	}
	cfg := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: llevel,
			},
		),
	}

	if driver == "sqlite" {
		DBConn, err = gorm.Open(sqlite.Open(database), cfg)
		if err != nil {
			fmt.Println("db err: ", err)
		}
	} else if driver == "postgres" {
		dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, username, database, password)
		DBConn, err = gorm.Open(postgres.Open(dsn), cfg)
		if err != nil {
			fmt.Println("db err: ", err)
		}
	} else if driver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)
		DBConn, err = gorm.Open(mysql.Open(dsn), cfg)
		if err != nil {
			fmt.Println("db err: ", err)
		}
	} else {
		DBConn, err = gorm.Open(sqlite.Open("simple-startpage.db"), &gorm.Config{})
		if err != nil {
			fmt.Println("db err: ", err)
		}
	}
}
