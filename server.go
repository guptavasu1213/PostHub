// Vasu Gupta
// ID: 3066521
// Assignment 1

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRetrievePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Retrieve single record")
}

// Retrieve all public posts
func handleRetrievePosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Retrieve public records")
}

// Delete the post with the given link if having edit access
func handleDeletePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete post")
}

// Update the post with the given link if having edit access
func handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update post")
}

// Create a post with the contents given by client and respond back with the access links
func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create Post")
}

func main() {
	r := mux.NewRouter()

	// For all posts
	r.Path("/api/v1/posts").Methods("GET").HandlerFunc(handleRetrievePosts)
	r.Path("/api/v1/posts").Methods("POST").HandlerFunc(handleCreatePost)

	// For individual post
	// ##########CHANGE IT TO PROPER REGEX FOR HEX
	r.Path("/api/v1/posts/{*}").Methods("GET").HandlerFunc(handleRetrievePost)
	r.Path("/api/v1/posts/{*}").Methods("UPDATE").HandlerFunc(handleUpdatePost)
	r.Path("/api/v1/posts/{*}").Methods("DELETE").HandlerFunc(handleDeletePost)

	log.Fatal(http.ListenAndServe(":8001", r))
}
