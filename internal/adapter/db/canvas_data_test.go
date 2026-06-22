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

	wb := &WhiteboardCanvasData{
		WhiteboardId: 1,
		StartX:       10,
		EndX:         20,
		StartY:       30,
		EndY:         40,
	}

	err := wb.Create(&core.CanvasData{
		WhiteboardId: 1,
		StartX:       10,
		EndX:         20,
		StartY:       30,
		EndY:         40,
	})
	assert.NoError(t, err)
	assert.NotZero(t, wb.ID)
	assert.Equal(t, uint(1), wb.WhiteboardId)
	assert.Equal(t, uint(10), wb.StartX)
	assert.Equal(t, uint(20), wb.EndX)
	assert.Equal(t, uint(30), wb.StartY)
	assert.Equal(t, uint(40), wb.EndY)
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

	model := &WhiteboardCanvasData{}
	got, err := model.GetByWhiteboardID(1)
	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, 1, got[0].WhiteboardId)
	assert.Equal(t, 10, got[0].StartX)
	assert.Equal(t, 12, got[1].StartX)
}
