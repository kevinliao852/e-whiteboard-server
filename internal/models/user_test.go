package models

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

func TestUser_Create(t *testing.T) {
	setupUserTestDB(t)

	repo := &User{}
	user := &core.User{
		DisplayName: "Test User",
		Email:       "test@example.com",
		GoogleID:    "google-123",
	}

	err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUser_GetById(t *testing.T) {
	setupUserTestDB(t)

	repo := &User{}
	user := &core.User{
		DisplayName: "Test User",
		Email:       "test@example.com",
		GoogleID:    "google-123",
	}
	_ = repo.Create(user)

	found, err := repo.GetById("1")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.DisplayName, found.DisplayName)
}

func TestUser_GetByGoogleId(t *testing.T) {
	setupUserTestDB(t)

	repo := &User{}
	user := &core.User{
		DisplayName: "Test User",
		Email:       "test@example.com",
		GoogleID:    "google-123",
	}
	_ = repo.Create(user)

	found, err := repo.GetByGoogleId("google-123")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.GoogleID, found.GoogleID)
}
