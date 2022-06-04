package main

import (
	"log"

	"github.com/zayyanmatin/sre-data-service/api/internal"
	"github.com/zayyanmatin/sre-data-service/pkg/sqlserver"
)

func main() {
	db, err := sqlserver.OpenDb()
	if err != nil {
		log.Fatal(err)
	}
	err = sqlserver.SetupDb(db)
	if err != nil {
		log.Fatal(err)
	}
	server, err := internal.Start(db)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	// Start server
	err = server.Api.Start(":1323")
	if err != nil {
		log.Fatal(err)
	}
}
