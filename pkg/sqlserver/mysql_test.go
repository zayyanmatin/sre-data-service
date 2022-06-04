package sqlserver

import (
	"database/sql/driver"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestSetup(t *testing.T) {
	// creating sql db and mock instance
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatal(err)
	}
	//mocking expected calls
	mock.ExpectPing()
	mock.ExpectExec("CREATE DATABASE").WillReturnResult(driver.ResultNoRows)
	mock.ExpectExec("CREATE TABLE").WillReturnResult(driver.ResultNoRows)
	//testing setup
	err = setupDb(db)
	if err != nil {
		t.Error(err)
	}
	// we make sure that all expectations were met
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
