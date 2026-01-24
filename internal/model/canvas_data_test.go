package model

import (
	"testing"

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

	err := Create(wb)
	assert.NoError(t, err)
	assert.NotZero(t, wb.ID)
	assert.Equal(t, uint(1), wb.WhiteboardId)
	assert.Equal(t, uint(10), wb.StartX)
	assert.Equal(t, uint(20), wb.EndX)
	assert.Equal(t, uint(30), wb.StartY)
	assert.Equal(t, uint(40), wb.EndY)
}
