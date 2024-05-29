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

	var gormDB *gorm.DB = &gorm.DB{}
	dialector := sqlite.Open("filename")

	m := mock_database.NewMockDBCreator(ctrl)

	// return a mock function
	m.EXPECT().CreateOpenGorm().Return(
		func(dialector gorm.Dialector, config *gorm.Config) (db *gorm.DB, err error) {
			return gormDB, nil
		},
	)
	m.EXPECT().CreateDialector().Return(dialector)

	if db, err := Connect(m); err != nil {
		t.Fatalf(err.Error())
	} else {
		t.Logf("Connected to database: %v", db)
	}

}
