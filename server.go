// Vasu Gupta
// ID: 3066521
// Assignment 2

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func main() {
	// Set up the Database
	var err error
	db, err = sqlx.Connect("sqlite3", "assign1.db")
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
	apiRouter.Path("/posts/{link_id:[0-9a-zA-Z]{32}}").Methods("POST").HandlerFunc(handleUpdatePost)
	apiRouter.Path("/posts/report/{link_id:[0-9a-zA-Z]{32}}").Methods("POST").HandlerFunc(handlePostReport)
	apiRouter.Path("/posts/{link_id:[0-9a-zA-Z]{32}}").Methods("DELETE").HandlerFunc(handleDeletePost)

	// Serve files
	r.PathPrefix("/posts/{link_id:[0-9a-zA-Z]{32}}").HandlerFunc(serveIndividualPostPage)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("dist")))

	portNumber := ":8010"

	fmt.Println("listening on port", portNumber)

	log.Fatal(http.ListenAndServe(portNumber, r))
}
