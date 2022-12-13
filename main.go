package main

import (
	"app/database"
	"app/models"
	"app/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func init() {

	sqlc := database.SQLiteCreate{
		OpenSQLiteFunc: sqlite.Open,
		OpenGormFunc:   gorm.Open,
	}

	database.DB, err = database.Connect(os.Getenv("DATABASE_PATH"), &sqlc)

	if err != nil {
		log.Fatal(err.Error())
	}

	database.DB.AutoMigrate(&models.User{})
	database.DB.AutoMigrate(&models.Whiteboard{})

	log.Print("Database is connected")
}

var err error

func main() {

	r := routes.Handler()
	r.Run(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))

}
