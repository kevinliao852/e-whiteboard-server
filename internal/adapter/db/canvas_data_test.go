package db

import (
	"testing"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
	"github.com/kevinliao852/e-whiteboard-server/internal/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupCanvasDataTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	database.DB = db

	err = database.DB.AutoMigrate(&WhiteboardCanvasData{})
	assert.NoError(t, err)
}

func TestCreate(t *testing.T) {
	setupCanvasDataTestDB(t)

	err := CreateCanvasData(&core.CanvasData{
		WhiteboardId: 1,
		StartX:       10,
		EndX:         20,
		StartY:       30,
		EndY:         40,
	})
	assert.NoError(t, err)

	var rows []WhiteboardCanvasData
	err = database.DB.Find(&rows).Error
	assert.NoError(t, err)
	assert.Len(t, rows, 1)
	assert.NotZero(t, rows[0].ID)
	assert.Equal(t, uint(1), rows[0].WhiteboardId)
	assert.Equal(t, uint(10), rows[0].StartX)
	assert.Equal(t, uint(20), rows[0].EndX)
	assert.Equal(t, uint(30), rows[0].StartY)
	assert.Equal(t, uint(40), rows[0].EndY)
}

func TestGetByWhiteboardID(t *testing.T) {
	setupCanvasDataTestDB(t)

	rows := []*WhiteboardCanvasData{
		{WhiteboardId: 1, StartX: 10, StartY: 20, EndX: 30, EndY: 40},
		{WhiteboardId: 2, StartX: 11, StartY: 21, EndX: 31, EndY: 41},
		{WhiteboardId: 1, StartX: 12, StartY: 22, EndX: 32, EndY: 42},
	}
	for _, row := range rows {
		err := database.DB.Create(row).Error
		assert.NoError(t, err)
	}

	got, err := GetCanvasDataByWhiteboardID(1)
	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, 1, got[0].WhiteboardId)
	assert.Equal(t, 10, got[0].StartX)
	assert.Equal(t, 12, got[1].StartX)
}
