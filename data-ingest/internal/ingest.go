package internal

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/zayyanmatin/sre-data-service/models"

	"github.com/zayyanmatin/sre-data-service/pkg/sqlserver"
)

type Server struct {
	db *sql.DB
}

func (s *Server) IngestData() error {
	// starting ingestion 5 minutes before current time
	before := time.Now().Add(-5 * time.Minute)
	// ending ingestion up until current time
	end := time.Now()
	for end.After(before) {
		ts := models.Timeseries{
			Timestamp:   uint32(before.Unix()),
			Cpu:         rand.Float32() * 100,
			Concurrency: uint32(rand.Intn(500000)),
		}
		// inserting data into timeseries database
		err := s.Insert(ts)
		if err != nil {
			return fmt.Errorf("could not insert datum: %w", err)
		}
		// incrementing to next instance of timeseries
		before = before.Add(60 * time.Second)
	}
	return nil
}

func (s *Server) Insert(datum models.Timeseries) error {
	// preparing sql execution statement
	stmt, err := s.db.Prepare("insert into sre.timeseries set ts=?, cpu=?, concurrency=?")
	if err != nil {
		return fmt.Errorf("could not prepare insert statement: %w", err)
	}
	defer stmt.Close()
	// executing statement with arguments
	_, err = stmt.Exec(datum.Timestamp, datum.Cpu, datum.Concurrency)
	if err != nil {
		return fmt.Errorf("could not execute insert statement: %w", err)
	}

	return nil
}

func Start() (*Server, error) {
	// connecting to mysql db
	db, err := sqlserver.OpenDb()
	if err != nil {
		return nil, fmt.Errorf("could not open mysql db: %w", err)
	}
	return &Server{db}, nil
}
