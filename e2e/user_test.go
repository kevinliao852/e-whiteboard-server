package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kevinliao852/e-whiteboard-server/internal/database"
	"github.com/kevinliao852/e-whiteboard-server/internal/models"
	"github.com/kevinliao852/e-whiteboard-server/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserCon(t *testing.T) {
	var err error
	gormConnector := database.NewGormConnector(os.Getenv("DATABASE_PATH"))

	database.DB, err = gormConnector.Connect()
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
