package db

import (
	"testing"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
	"github.com/kevinliao852/e-whiteboard-server/internal/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	database.DB = db

	err = database.DB.AutoMigrate(&User{})
	assert.NoError(t, err)
}

func TestUserModel_CreateAndGet(t *testing.T) {
	setupUserTestDB(t)

	model := &User{}
	err := model.Create(&core.User{
		DisplayName: "Test User",
		Email:       "test@example.com",
		GoogleID:    "google-1",
	})
	assert.NoError(t, err)

	user, err := model.GetByGoogleId("google-1")
	assert.NoError(t, err)
	assert.Equal(t, "Test User", user.DisplayName)
	assert.Equal(t, "test@example.com", user.Email)
}
