package main

import (
	"log"

	"github.com/zayyanmatin/sre-data-service/api/internal"
	"github.com/zayyanmatin/sre-data-service/pkg/sqlserver"
)

func main() {
	//fetching db
	db := sqlserver.FetchDb()
	server, err := internal.Setup(db)
	if err != nil {
		log.Fatal("unable to setup server", err)
	}
	// closing resources before exit
	defer func() {
		if err = server.Close(); err != nil {
			log.Fatal("unable to close server", err)
		}
	}()

	// Start server
	err = server.Api.Start(":1323")
	if err != nil {
		log.Fatal("unable to start api", err)
	}
}
