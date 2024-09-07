package main

import (
	"app/database"
	"app/models"
	"app/pkg/config"
	"app/routes"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configManager := config.NewConfigManager([]string{"APP_HOST", "APP_PORT"})

	err = configManager.CheckAndLoadConfig()

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.Printf("Log level is %v", log.GetLevel())
	log.Printf("Set log format to JSON")
}

func init() {
	sqlc := database.SQLiteCreate{
		OpenSQLiteFunc: sqlite.Open,
		OpenGormFunc:   gorm.Open,
		Filename:       os.Getenv("DATABASE_PATH"),
	}

	database.DB, err = database.Connect(&sqlc)

	if err != nil {
		log.Fatal(err.Error())
	}

	err = database.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = database.DB.AutoMigrate(&models.Whiteboard{})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = database.DB.AutoMigrate(&models.WhiteboardCanvasData{})
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Print("Database is connected")
}

var err error

func main() {
	wc := routes.WithCORS()

	r := routes.Handler(wc)
	r.Run(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))

}
