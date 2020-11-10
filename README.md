# PostHub

## Overview
This repository contains a RESTful API and front-end for a text sharing web application where the users can share their posts with the public while maintaining anonymity. 

## Database used
* Sqlite 3

## Prerequisites
* Golang (Tested on 1.15.2)
Follow these steps: https://golang.org/doc/install
* Gorrila Mux Package
```
go get -u github.com/gorilla/mux
```
* Sqlx Library
```
go get github.com/jmoiron/sqlx
```
* Goose database migration tool
```
go get -u github.com/pressly/goose/cmd/goose
```
* Sqlite3 driver for golang
```
go get github.com/mattn/go-sqlite3
``` 
* Typescript (Tested on 3.8.2)
```
sudo apt install node-typescript
```
* Sqlite3 for local use [optional] (Tested on 3.31.1)
```
sudo apt install sqlite3
```

## Usage
To build the server, run the following:
* -port is an optional argument that specifies the port number to launch the server on. By default, the port is 8080.
* -dbPath is an optional argument that specifies the database file to import and export the records from. By default, it uses `records.db` file
```
go build -o server
./server [-port <port number>] [-dbPath <database file path>]
```
After running the server, open your web browser and type this in the address bar: `localhost:<port number of go server>`

## Example
```
./server -port 8080 -dbPath ./someDatabaseFile.db
```

## Status
Completed

## Authors
* Vasu Gupta