package main

import (
	"app/routes"
	"app/database"
	"app/models"
)

var err error

func main() {
	database.DB, err = database.Connect()

	database.DB.AutoMigrate(&models.User{})

	r := routes.Handler()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
