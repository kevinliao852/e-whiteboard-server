package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"app/database"
	"app/models"
	"app/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserCon(t *testing.T) {
	sqlc := database.SQLiteCreate{
		OpenSQLiteFunc: sqlite.Open,
		OpenGormFunc:   gorm.Open,
		Filename:       os.Getenv("DATABASE_PATH"),
	}

	var err error

	database.DB, err = database.Connect(&sqlc)

	if err != nil {
		panic("failed to connect database")
	}

	err = database.DB.AutoMigrate(&models.User{})

	if err != nil {
		panic("failed to connect database")
	}

	assert.Nil(t, err)

	cam := routes.WithAuthMiddleware(func(c *gin.Context) {
		c.Next()
	})

	router := routes.Handler(cam)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/user", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "[]", w.Body.String())

}
