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
	err = database.DB.AutoMigrate(&WhiteboardCanvasData{})
	assert.NoError(t, err)
}

func TestWhiteboardModel_CreateAndGetByUserId(t *testing.T) {
	setupWhiteboardTestDB(t)

	err := CreateWhiteboard(&core.Whiteboard{
		UserId: 1,
		Name:   "Board 1",
	})
	assert.NoError(t, err)

	boards, err := GetWhiteboardsByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, boards, 1)
	assert.Equal(t, "Board 1", boards[0].Name)
}

func TestWhiteboardModel_DeleteRemovesCanvasData(t *testing.T) {
	setupWhiteboardTestDB(t)

	err := database.DB.Create(&Whiteboard{
		Id:   10,
		Name: "Board 10",
	}).Error
	assert.NoError(t, err)

	err = database.DB.Create(&WhiteboardCanvasData{
		WhiteboardId: 10,
		StartX:       1,
		StartY:       2,
		EndX:         3,
		EndY:         4,
	}).Error
	assert.NoError(t, err)

	err = DeleteWhiteboard(10)
	assert.NoError(t, err)

	var whiteboards []Whiteboard
	err = database.DB.Find(&whiteboards).Error
	assert.NoError(t, err)
	assert.Len(t, whiteboards, 0)

	var canvasRows []WhiteboardCanvasData
	err = database.DB.Find(&canvasRows).Error
	assert.NoError(t, err)
	assert.Len(t, canvasRows, 0)
}
