package database

import (
	"app/mock/mock_database"
	"testing"

	"github.com/golang/mock/gomock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestConnect(t *testing.T) {
	ctrl := gomock.NewController(t)

	var db *gorm.DB = &gorm.DB{}
	dialector := sqlite.Open("filename")

	m := mock_database.NewMockSQLiteCreator(ctrl)

	m.EXPECT().OpenGorm(dialector, &gorm.Config{}).Return(db, nil)
	m.EXPECT().OpenSQLite(gomock.Any()).Return(dialector)

	if _, err := Connect("filename", m); err != nil {
		t.Fatalf(err.Error())
	}

}
