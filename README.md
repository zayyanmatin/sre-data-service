# SRE Data Service
A Service run in Go that will ingest mock cpu and concurrency data into mysql.
A REST Api is provided to query the data.

## Prerequisites and Setup
1) Have MySQL installed and start up on a server.
2) Create a DataSourceName (DSN) for the database, making sure the user has privileges to insert and query data into a database called 'sre'.
3) Once your DSN is ready, set it to the SRE_DSN environment variable e.g:
export SRE_DSN=your_dsn

## Ingest Data
To ingest data into the database, run main.go from data-ingest/cmd

## API
To start the API server, run main.go from api/cmd. API will be available on port 1323

API specification is available within the repository

## Running Unit Tests
From repository's root directory run go test ./... -v

## Enhancements to be made
Add validation on startTime and endTime with 400 bad requests 

Add logging to both ingestion and api 

Make updates and corrections to API spec

## Tech debt
Add script to start sql server and securely pass DSN before running main.go file

Update API Unit test to check response body

Increase test coverage - check unhappy paths

Add optional flag to clean up and drop database

Add configs file to manage environment variables properly
