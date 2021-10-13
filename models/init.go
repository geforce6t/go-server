package models

import (
	"log"

	"github.com/geforce6t/go-server/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := utils.GetEnvValue("user") + ":" + utils.GetEnvValue("pass") + "@tcp(127.0.0.1:3306)/" + utils.GetEnvValue("dbname") +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error establishing a connection to the database!: %v", err)
	}

	log.Println("Connected to the database", utils.GetEnvValue("dbname"))

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Error while migration: %v", err)
	}

	return db
}
