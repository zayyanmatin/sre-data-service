package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/zayyanmatin/sre-data-service/data-ingest/internal"
)

func init() {
	// initialize random data source
	rand.Seed(time.Now().UnixNano())
}

func main() {
	//starting db server
	server, err := internal.Start()
	if err != nil {
		log.Fatal(err)
	}
	//starting ingestion
	if err := server.IngestData(); err != nil {
		log.Fatal(err)
	}
}
