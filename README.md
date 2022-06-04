# sre-data-service
A Service run in Go that will ingest mock cpu and concurrency data into mysql.
A REST Api is provided to query the data.

## Prerequisites and Setup
1) Have MySQL installed and start up on a server.
2) Create a DataSourceName (DSN) for the database, making sure the user has privileges to insert and query data into a database called 'sre'.
3) Once a DSN is ready, set it to the SRE_DSN environment variable e.g:
export SRE_DSN=your_dsn

## Ingest Data
To ingest data into the database, run main.go from data-ingest/cmd

## API
To start the API server, run main.go from api/cmd. API will be available on port 1323.

## Enhancements to be made
validation on startTime and endTime with 400 bad requests 
make startTime and endTime time variables instead of integers
updates to description in API spec

## Tech debt
add script to start sql server before running main.go file

add optional flag to clean up and drop database
add configs file to manage environment variables properly
