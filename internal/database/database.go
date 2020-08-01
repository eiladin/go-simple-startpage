package database

import (
	"fmt"
	"strings"

	"github.com/eiladin/go-simple-startpage/internal/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres
	_ "github.com/jinzhu/gorm/dialects/sqlite"   //sqlite
)

// DBConn holds the database connection
var DBConn *gorm.DB

// InitDB initialized the selected database
func InitDB() *gorm.DB {
	var d *gorm.DB
	var err error
	c := config.GetConfig()

	driver := strings.ToLower(c.Database.Driver)
	database := c.Database.Dbname
	username := c.Database.Username
	password := c.Database.Password
	host := c.Database.Host
	port := c.Database.Port
	log := c.Database.Log

	if driver == "sqlite" {
		d, err = gorm.Open("sqlite3", database)
		if err != nil {
			fmt.Println("db err: ", err)
		}
	} else if driver == "postgres" {
		d, err = gorm.Open("postgres", "host="+host+" port="+port+" user="+username+" dbname="+database+"  sslmode=disable password="+password)
		if err != nil {
			fmt.Println("db err: ", err)
		}
	} else if driver == "mysql" {
		d, err = gorm.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/"+database+"?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			fmt.Println("db err: ", err)
		}
	} else {
		d, err = gorm.Open("sqlite3", "simple-startpage.db")
		if err != nil {
			fmt.Println("db err: ", err)
		}
	}

	d.LogMode(log)
	DBConn = d
	return DBConn
}
