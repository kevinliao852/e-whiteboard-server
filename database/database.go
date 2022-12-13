package database

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

type SQLiteCreator interface {
	OpenGorm(dialector gorm.Dialector, config *gorm.Config) (db *gorm.DB, err error)
	OpenSQLite(dsn string) gorm.Dialector
}

type SQLiteCreate struct {
	OpenSQLiteFunc func(dsn string) gorm.Dialector
	OpenGormFunc   func(dialector gorm.Dialector, config *gorm.Config) (db *gorm.DB, err error)
}

func (sqlc *SQLiteCreate) OpenGorm(dialector gorm.Dialector, config *gorm.Config) (db *gorm.DB, err error) {
	return sqlc.OpenGormFunc(dialector, config)
}

func (sqlc *SQLiteCreate) OpenSQLite(filename string) gorm.Dialector {
	return sqlc.OpenSQLiteFunc(filename)
}

func Connect(filename string, createSQLiter SQLiteCreator) (*gorm.DB, error) {
	db, err := createSQLiter.OpenGorm(createSQLiter.OpenSQLite(filename), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db, err
}
