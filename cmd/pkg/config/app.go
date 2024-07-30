package config

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// this file will export the db to other files such as the models and the controllers
var (
	db *gorm.DB
)

func Connect() {
	d, err := gorm.Open("mysql", "root:soumil05/goproj?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	db = d
}

func GetDB() *gorm.DB {
	return db
}
