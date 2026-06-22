package main

import (
	"os"

	"github.com/kevinliao852/e-whiteboard-server/internal/adapter/db"
	"github.com/kevinliao852/e-whiteboard-server/internal/database"
	"github.com/kevinliao852/e-whiteboard-server/internal/route"
	"github.com/kevinliao852/e-whiteboard-server/pkg/config"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// Load environment variables and check required config
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	configManager := config.NewConfigManager([]string{
		"APP_HOST",
		"APP_PORT",
		"SESSION_SECRET",
	})

	if err = configManager.CheckAndLoadConfig(); err != nil {
		log.Fatal(err)
	}
}

// Set log level and format
func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.Printf("Log level is %v", log.GetLevel())
	log.Printf("Set log format to JSON")
}

// Connect to the database and run migrations
func init() {
	gormConector := database.NewGormConnector(os.Getenv("DATABASE_PATH"))
	database.DB, err = gormConector.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	err = database.DB.AutoMigrate(&db.User{})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = database.DB.AutoMigrate(&db.Whiteboard{})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = database.DB.AutoMigrate(&db.WhiteboardCanvasData{})
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Print("Database is connected")
}

var err error

func main() {
	wc := route.WithCORS()

	r := route.Handler(wc)
	err := r.Run(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))

	if err != nil {
		log.Fatal(err)
	}
}
