package main

import (
	"app/database"
	"app/models"
	"app/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var err error

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.DB, _ = database.Connect(os.Getenv("DATABASE_PATH"))
	database.DB.AutoMigrate(&models.User{})
	database.DB.AutoMigrate(&models.Whiteboard{})

	r := routes.Handler()
	r.Run(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))

}
