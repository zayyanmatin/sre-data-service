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
	//fetching db
	db := sqlserver.FetchDb()
	server := internal.Server{
		Db: db,
	}
	//starting ingestion
	if err := server.IngestData(); err != nil {
		log.Fatal(err)
	}
}
