package database

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

type Database interface {
}

type DBCreator interface {
	CreateDialector() gorm.Dialector
	CreateOpenGorm() func(dialector gorm.Dialector, config *gorm.Config) (db *gorm.DB, err error)
}

type SQLiteCreate struct {
	OpenSQLiteFunc func(dsn string) gorm.Dialector
	OpenGormFunc   func(dialector gorm.Dialector, config *gorm.Config) (db *gorm.DB, err error)
	Filename       string
}

func (sqlc *SQLiteCreate) CreateDialector() gorm.Dialector {
	return sqlc.OpenSQLiteFunc(sqlc.Filename)
}

func (sqlc *SQLiteCreate) CreateOpenGorm() func(dialector gorm.Dialector, config *gorm.Config) (db *gorm.DB, err error) {
	return sqlc.OpenGormFunc
}

func Connect(creator DBCreator) (*gorm.DB, error) {
	db, err := creator.CreateOpenGorm()(creator.CreateDialector(), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	return db, err
}
