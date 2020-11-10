package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var portNumber int
var databaseFilePath string

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options]\n\nOptions:\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
}

func parseFlags() {
	flag.Usage = usage
	flag.IntVar(&portNumber, "port", 8080, "port number for connection")
	flag.StringVar(&databaseFilePath, "dbPath", "records.db", "database file path")
	flag.Parse()
}

func main() {
	parseFlags()

	// Set up the Database
	var err error
	db, err = sqlx.Connect("sqlite3", databaseFilePath)
	defer db.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot connect to database: %v\n", err)
		os.Exit(1)
	}

	// Set up router and subrouter
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	// For all posts
	apiRouter.Path("/posts").Methods("GET").HandlerFunc(handleRetrievePosts)
	apiRouter.Path("/posts").Methods("POST").HandlerFunc(handleCreatePost)

	// For an individual post
	apiRouter.Path("/posts/{link_id:[0-9a-zA-Z]{32}}").Methods("GET").HandlerFunc(handleRetrievePost)
	apiRouter.Path("/posts/{link_id:[0-9a-zA-Z]{32}}").Methods("PUT").HandlerFunc(handleUpdatePost)
	apiRouter.Path("/posts/{link_id:[0-9a-zA-Z]{32}}/report").Methods("POST").HandlerFunc(handlePostReport)
	apiRouter.Path("/posts/{link_id:[0-9a-zA-Z]{32}}").Methods("DELETE").HandlerFunc(handleDeletePost)

	// Serve files
	r.Path("/pastes/{link_id:[0-9a-zA-Z]{32}}").HandlerFunc(handleIndividualPageServing)
	r.Path("/pastes").HandlerFunc(handlerToRetrieveAllPostsPage)
	r.Path("/").HandlerFunc(handlerToRetrieveHomePage)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("dist")))

	fmt.Printf("listening on port :%d\n", portNumber)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNumber), r))
}
