package models

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetEnvValue(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func InitDB() *gorm.DB {
	dsn := GetEnvValue("user") + ":" + GetEnvValue("pass") + "@tcp(127.0.0.1:3306)/" + GetEnvValue("dbname") +
		"?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error establishing a connection to the database!: %v", err)
	}

	log.Println("Connected to the database", GetEnvValue("dbname"))

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Error while migration: %v", err)
	}

	return db
}
