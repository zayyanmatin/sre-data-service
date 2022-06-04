package sqlserver

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("SRE_DSN"))
	if err != nil {
		return nil, fmt.Errorf("could not open mysql: %w", err)
	}
	return db, nil
}

func SetupDb(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("could not ping mysql: %w", err)
	}
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS sre;")
	if err != nil {
		return fmt.Errorf("could not execute query for database creation: %w", err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS sre.timeseries (ts integer, cpu float, concurrency integer);")
	if err != nil {
		return fmt.Errorf("could not execute query for table creation: %w", err)
	}
	fmt.Println("successfully connected to mysql")
	return nil
}