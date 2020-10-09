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
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type postJSON struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Scope string `json:"scope"`
}

type postLinks struct {
	EditLink string `json:"editLink,omitempty"`
	ViewLink string `json:"viewLink"`
}

// Used when retrieving posts from the database
type postRetrieval struct {
	PostID int64  `json:"-" db:"post_id,omitempty"`
	Title  string `json:"title,omitempty" db:"title,omitempty"`
	Body   string `json:"body,omitempty" db:"body,omitempty"`
	Scope  string `json:"scope,omitempty" db:"scope,omitempty"` // Private or Public
	LinkID string `json:"-" db:"link_id,omitempty"`
	Access string `json:"-" db:"access,omitempty"` // Edit or View
	Epoch  string `json:"epoch,omitempty" db:"epoch,omitempty"`
}

var db *sqlx.DB

// ADD A MUTEX?########

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

	lastIndex := strings.LastIndex(r.URL.Path, "/")
	userLinkID := r.URL.Path[lastIndex+1:]
	resourceWithoutLinkID := r.URL.Path[:lastIndex]

	// Retrieve the post from the database
	query := `SELECT *
			 	FROM Posts, Links 
			 	WHERE Posts.post_id = Links.post_id and link_id = $1`
	entry := postRetrieval{}
	err := db.Get(&entry, query, userLinkID)
	if err == sql.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		fmt.Println("no Entries found")
	} else if err != nil {
		fmt.Println("unsuccessful data lookup", err)
	}

	if entry.Access == "View" {
		entry.Scope = "" // To avoid JSON encoding of the scope

		// Encode and Send Response To Client
		err = json.NewEncoder(w).Encode(entry)
		if err != nil {
			log.Print("error: encoding unsuccessful")
		}

	} else { // Edit Access
		// look up the view-id
		var viewID string

		query = `SELECT link_id 
					FROM Links
					WHERE post_id = $1 and access = "View"`
		err = db.Get(&viewID, query, entry.PostID)
		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			fmt.Println("no Entries found")
		} else if err != nil {
			fmt.Println("unsuccessful data lookup", err)
		}

		// Add Links to struct
		links := postLinks{
			EditLink: path.Join(r.Host, resourceWithoutLinkID, entry.LinkID),
			ViewLink: path.Join(r.Host, resourceWithoutLinkID, viewID),
		}

		// Combine the post and links structs for JSON encoding
		response := struct {
			postRetrieval
			postLinks
		}{entry, links}

		// Encode and Send Response To Client
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Print("error: encoding unsuccessful")
		}

	}
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
	updatedPost := postJSON{}
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

// isUniqueViolation returns true if the supplied error resulted from a primary key constraint failure.
func isUniqueViolation(err error) bool {
	if err, ok := err.(sqlite3.Error); ok {
		return err.Code == 19 && err.ExtendedCode == 1555
	}

	return false
}

// Return the unique link id and add it to the database
func addLinkIDToDatabase(linkID string, postID int64, access string) string {
	query := `INSERT INTO links (link_id, access, post_id)
                       VALUES ($1, $2, $3)`
	_, err := db.Exec(query, linkID, access, postID)
	if err != nil {
		if isUniqueViolation(err) {
			linkID = addLinkIDToDatabase(generateRandomString(), postID, access)
		} else {
			fmt.Errorf("post creation unsuccessful")
		}
	}
	return linkID
}

// Create a post with the contents given by client and respond back with the access links
func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create")
	// Decode Post Contents
	newPost := postJSON{}
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		log.Println("error: decoding error occured")
	}

	fmt.Println(newPost.Title)

	fmt.Println(r.Host)

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
	editID := addLinkIDToDatabase(generateRandomString(), postID, "Edit")
	viewID := addLinkIDToDatabase(generateRandomString(), postID, "View")

	editLink := getLinkFromID(r, editID)
	viewLink := getLinkFromID(r, viewID)
	fmt.Println(editLink, viewLink)

	// Encode and Send Response To Client
	response := postLinks{EditLink: editLink, ViewLink: viewLink}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Print("error: encoding unsuccessful")
	}
}

func main() {
	var err error
	db, err = sqlx.Connect("sqlite3", "assign1.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot connect to database: %v\n", err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	// For all posts
	// r.Path("/api/v1/posts").Methods("GET").HandlerFunc(handleRetrievePosts)
	r.Path("/api/v1/posts").Methods("POST").HandlerFunc(handleCreatePost)

	// For individual post
	// ##########CHANGE IT TO PROPER REGEX FOR HEX
	r.Path("/api/v1/posts/{*}").Methods("GET").HandlerFunc(handleRetrievePost)
	r.Path("/api/v1/posts/{*}").Methods("UPDATE").HandlerFunc(handleUpdatePost)
	r.Path("/api/v1/posts/{*}").Methods("DELETE").HandlerFunc(handleDeletePost)

	log.Fatal(http.ListenAndServe(":8101", r))
}
