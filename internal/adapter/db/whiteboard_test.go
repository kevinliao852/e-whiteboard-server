package db

import (
	"testing"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
	"github.com/kevinliao852/e-whiteboard-server/internal/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupWhiteboardTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	database.DB = db

	err = database.DB.AutoMigrate(&Whiteboard{})
	assert.NoError(t, err)
}

func TestWhiteboardModel_CreateAndGetByUserId(t *testing.T) {
	setupWhiteboardTestDB(t)

	model := &Whiteboard{}
	err := model.Create(&core.Whiteboard{
		UserId: 1,
		Name:   "Board 1",
	})
	assert.NoError(t, err)

	boards, err := model.GetByUserId(1)
	assert.NoError(t, err)
	assert.Len(t, boards, 1)
	assert.Equal(t, "Board 1", boards[0].Name)
}
