package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type DBConnector interface {
	Connect() (*gorm.DB, error)
}

var _ DBConnector = &GormConnector{}

type GormConnector struct {
	FileName string
}

func NewGormConnector(fileName string) *GormConnector {
	return &GormConnector{FileName: fileName}
}

func (gc *GormConnector) Connect() (*gorm.DB, error) {
	dialector := sqlite.Open(gc.FileName)
	db, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	return db, err
}
