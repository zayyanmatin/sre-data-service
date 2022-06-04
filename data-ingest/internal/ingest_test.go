package internal

import (
	"database/sql/driver"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestIngestion(t *testing.T) {
	// creating sql db and mock instance
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	// creating server instance with mock db
	server := Server{
		db: db,
	}
	// sql statements prepared and executed every minute for 5 minutes per ingestion
	for i := 0; i < 6; i++ {
		mock.ExpectPrepare("insert into").WillBeClosed().ExpectExec().WillReturnResult(driver.ResultNoRows)
	}
	// starting ingestion
	err = server.IngestData()
	if err != nil {
		t.Error(err)
	}
	// ensuring expectations were met
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
