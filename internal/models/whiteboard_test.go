package models

import (
	"testing"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
	"github.com/kevinliao852/e-whiteboard-server/internal/database"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB initializes an in-memory SQLite database for tests
func setupTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// override global DB
	database.DB = db

	// migrate schema
	err = database.DB.AutoMigrate(&core.Whiteboard{})
	assert.NoError(t, err)
}

func TestWhiteboard_Create(t *testing.T) {
	setupTestDB(t)

	repo := &Whiteboard{}

	wb := &core.Whiteboard{
		UserId: 1,
		Name:   "Test Board",
	}

	err := repo.Create(wb)
	assert.NoError(t, err)
	assert.NotZero(t, wb.Id)
}

func TestWhiteboard_GetByUserId(t *testing.T) {
	setupTestDB(t)

	repo := &Whiteboard{}

	// seed data
	_ = repo.Create(&core.Whiteboard{UserId: 1, Name: "Board A"})
	_ = repo.Create(&core.Whiteboard{UserId: 1, Name: "Board B"})
	_ = repo.Create(&core.Whiteboard{UserId: 2, Name: "Board C"})

	wbs, err := repo.GetByUserId(1)
	assert.NoError(t, err)
	assert.Len(t, wbs, 2)
}

func TestWhiteboard_Delete(t *testing.T) {
	setupTestDB(t)

	repo := &Whiteboard{}

	wb := &core.Whiteboard{
		UserId: 1,
		Name:   "To Delete",
	}

	_ = repo.Create(wb)

	err := repo.Delete(wb.Id)
	assert.NoError(t, err)

	var count int64
	database.DB.Model(&core.Whiteboard{}).
		Where("id = ?", wb.Id).
		Count(&count)

	assert.Equal(t, int64(0), count)
}
