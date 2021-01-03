package database

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

type Config struct {
	filename string
}

var DB *gorm.DB

func GetConfig() *Config {
	config := Config {
		
	}
	return &config
}

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db, err
}
