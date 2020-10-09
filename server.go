// Vasu Gupta
// ID: 3066521
// Assignment 1

package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
)

type post struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Scope string `json:"scope"`
}

type postCreationResponse struct {
	EditLink string `json:"editLink"`
	ViewLink string `json:"viewLink"`
}

// Generate a random string of 32 characters
func generateRandomString() string {
	randomNum := strconv.Itoa(rand.Int())
	hash := md5.Sum([]byte(randomNum))
	return hex.EncodeToString(hash[:])
}

func handleRetrievePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Retrieve single record")
	// Extract the record corresponding to if the ID is an edit OR view
	// If not found any records, then send 404 or something
	log.Println(r.URL.RequestURI())

	_ = "SELECT * FROM Posts, Links WHERE Posts.post_id = Links.post_id and (edit_id = ### or view_id = ###)"
	// Check if the resource has the edit id
	// 		-> Then return title, body, scope, admin link, view link
	// Check if the resource has the view id
	// 		-> Then return title, body, scope
	// Else return 404 or SOMETHING (BAD REQUEST MAYBE?)#######

}

// Retrieve all public posts
func handleRetrievePosts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Retrieve public")
	// Extract the public records which aren't reported
	// If not found any records, then send nothing

	// USE PAGINATION-- offset stuff
	// Check boolean values in sqlite3
	_ = "SELECT * FROM Posts, Links, Report WHERE Posts.post_id = Links.post_id and Report.post_id = Posts.post_id and scope = \"Public\" and reported = \"False\""

	log.Println(r.URL.RequestURI())
}

// Delete the post with the given link if having edit access
func handleDeletePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Del")
	// DELETE the record corresponding to if the ID is an edit
	// If the link is view, send 403
	// If not found any records, then send 404 or something

	_ = "DELETE FROM Posts, Links WHERE Posts.post_id = Links.post_id and edit_id = ###"
	// Check if the resource has the edit id
	// 		-> DELETE the record corresponding to if the ID
	// Check if the resource has the view id
	// 		-> Send 403
	// Else return 404

	log.Println(r.URL.RequestURI())

}

// Update the post with the given link if having edit access
func handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update")
	// UPDATE the record corresponding to if the ID is an edit
	// If the link is view, send 403
	// If not found any records, then send 404 or something

	log.Println(r.URL.RequestURI())
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Decode Post Contents
	updatedPost := post{}
	err = json.Unmarshal(reqBody, &updatedPost)
	if err != nil {
		log.Println("error: decoding error occured")
	}
	fmt.Println(updatedPost)

	// Try updating and see if it updates the value not passed or not
	_ = "UPDATE Posts, Links SET ...... WHERE Posts.post_id = Links.post_id and edit_id = ###"
}

// Create a post with the contents given by client and respond back with the access links
func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Decode Post Contents
	newPost := post{}
	err = json.Unmarshal(reqBody, &newPost)
	if err != nil {
		log.Println("error: decoding error occured")
	}

	fmt.Println(newPost.Title)
	editLink := path.Join(r.RequestURI, generateRandomString())
	viewLink := path.Join(r.RequestURI, generateRandomString())
	fmt.Println(editLink, viewLink)

	// MIGHT HAVE TO CHANGE TO MARSHALL FOR CODE CONSISTENCY
	// Encode and Send Response To Client
	response := postCreationResponse{EditLink: editLink, ViewLink: viewLink}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Print("error: encoding unsuccessful")
	}

	// Store Values in the DB
}

func main() {
	r := mux.NewRouter()

	// For all posts
	// r.Path("/api/v1/posts").Methods("GET").HandlerFunc(handleRetrievePosts)
	r.Path("/api/v1/posts").Methods("POST").HandlerFunc(handleCreatePost)

	// For individual post
	// ##########CHANGE IT TO PROPER REGEX FOR HEX
	r.Path("/api/v1/posts/{*}").Methods("GET").HandlerFunc(handleRetrievePost)
	r.Path("/api/v1/posts/{*}").Methods("UPDATE").HandlerFunc(handleUpdatePost)
	r.Path("/api/v1/posts/{*}").Methods("DELETE").HandlerFunc(handleDeletePost)

	log.Fatal(http.ListenAndServe(":8001", r))
}
