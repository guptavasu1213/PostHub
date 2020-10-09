// Vasu Gupta
// ID: 3066521
// Assignment 1

package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	sqlite3 "github.com/mattn/go-sqlite3"
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

var db *sqlx.DB

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
	// Decode Post Contents
	updatedPost := post{}
	err := json.NewDecoder(r.Body).Decode(&updatedPost)

	if err != nil {
		log.Println("error: decoding error occured")
	}

	// Encode and Send Response To Client
	err = json.NewEncoder(w).Encode(updatedPost)
	if err != nil {
		log.Print("error: encoding unsuccessful")
	}

	fmt.Println(updatedPost)

	// Try updating and see if it updates the value not passed or not
	_ = "UPDATE Posts, Links SET ...... WHERE Posts.post_id = Links.post_id and edit_id = ###"
}

// isUniqueViolation returns true if the supplied error resulted from a unique constraint violation.
func isUniqueViolation(err error) bool {
	if err, ok := err.(sqlite3.Error); ok {
		return err.Code == 19 && err.ExtendedCode == 2067
	}

	return false
}

func addLinkIDToDatabase(linkID string, postID int64, access string) {
	query := `INSERT INTO links (link_id, access, post_id)
                       VALUES ($1, $2, $3)`
	_, err := db.Exec(query, linkID, access, postID)
	if err != nil {
		if isUniqueViolation(err) {
			addLinkIDToDatabase(linkID, postID, access)
		} else {
			fmt.Errorf("post creation unsuccessful")
		}
	}
}

// Create a post with the contents given by client and respond back with the access links
func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create")

	// Decode Post Contents
	newPost := post{}
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		log.Println("error: decoding error occured")
	}

	fmt.Println(newPost.Title)

	fmt.Println(r.Host)

	editID := generateRandomString()
	viewID := generateRandomString()

	editLink := path.Join(r.Host, r.RequestURI, editID)
	viewLink := path.Join(r.Host, r.RequestURI, viewID)
	fmt.Println(editLink, viewLink)

	// Insert data in posts table
	var result sql.Result
	query := `INSERT INTO posts (title, body, scope, epoch)
                       VALUES ($1, $2, $3, $4)`
	result, err = db.Exec(query, newPost.Title, newPost.Body, newPost.Scope, time.Now().Unix())
	if err != nil {
		fmt.Errorf("post creation unsuccessful")
	}
	postID, err := result.LastInsertId()

	// Insert the Link ID's to the Link table
	addLinkIDToDatabase(editID, postID, "Edit")
	addLinkIDToDatabase(viewID, postID, "View")

	// Encode and Send Response To Client
	response := postCreationResponse{EditLink: editLink, ViewLink: viewLink}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Print("error: encoding unsuccessful")
	}
}

func main() {
	r := mux.NewRouter()

	var err error
	db, err = sqlx.Connect("sqlite3", "assign1.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot connect to database: %v\n", err)
		os.Exit(1)
	}
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
