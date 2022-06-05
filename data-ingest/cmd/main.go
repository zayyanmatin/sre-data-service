package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/zayyanmatin/sre-data-service/data-ingest/internal"
	"github.com/zayyanmatin/sre-data-service/pkg/sqlserver"
)

func init() {
	// initialize random data source
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// fetching db
	db := sqlserver.FetchDb()
	// closing db resource before exit
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("unable to close db", err)
		}
	}()
	// create server instance with db
	server := internal.Server{
		Db: db,
	}
	// starting ingestion
	if err := server.IngestData(); err != nil {
		log.Fatal("error in ingesting data", err)
	}
}
