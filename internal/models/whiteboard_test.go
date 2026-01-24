package models

import (
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/kevinliao852/e-whiteboard-server/internal/database"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	now = time.Now()
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func setup(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
	}

	return mock, gormDB
}

func TestCreateAWhiteboard(t *testing.T) {
	mock, gormDB := setup(t)
	database.DB = gormDB

	wb := &Whiteboard{
		UserId:    1,
		Name:      "Test Whiteboard",
		CreatedAt: now,
		UpdatedAt: now,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `whiteboards` (`user_id`,`name`,`created_at`,`updated_at`) VALUES (?,?,?,?)")).
		WithArgs(wb.UserId, wb.Name, AnyTime{}, AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	var w Whiteboard
	if err := w.CreateAWhiteboard(wb); err != nil {
		t.Errorf("Failed to create whiteboard: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
