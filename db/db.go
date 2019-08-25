package db

import (
	"fmt"

	"github.com/eiladin/go-simple-startpage/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

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
		db, err = gorm.Open("sqlite3", "simple-startpage.db")
		if err != nil {
			fmt.Println("db err: ", err)
		}
	}

	db.LogMode(true)
	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
