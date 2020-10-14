// Vasu Gupta
// ID: 3066521
// Assignment 1

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
	apiRouter.Path("/posts/{link_id:[0-9a-zA-Z]+}").Methods("GET").HandlerFunc(handleRetrievePost)
	apiRouter.Path("/posts/{link_id:[0-9a-zA-Z]+}").Methods("POST").HandlerFunc(handleUpdatePost)
	apiRouter.Path("/posts/report/{link_id:[0-9a-zA-Z]+}").Methods("POST").HandlerFunc(handlePostReport)
	apiRouter.Path("/posts/{link_id:[0-9a-zA-Z]+}").Methods("DELETE").HandlerFunc(handleDeletePost)

	log.Fatal(http.ListenAndServe(":8010", r))
}
