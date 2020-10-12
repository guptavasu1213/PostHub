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

type postLinks struct {
	EditLink string `json:"editLink,omitempty"`
	ViewLink string `json:"viewLink"`
}

// Used when retrieving posts from the database
type post struct {
	PostID       int64  `json:"-" db:"post_id,omitempty"`
	Title        string `json:"title,omitempty" db:"title,omitempty"`
	Body         string `json:"body,omitempty" db:"body,omitempty"`
	Scope        string `json:"scope,omitempty" db:"scope,omitempty"` // Private or Public
	LinkID       string `json:"link,omitempty" db:"link_id,omitempty"`
	Access       string `json:"-" db:"access,omitempty"` // Edit or View
	Epoch        int64  `json:"epoch,omitempty" db:"epoch,omitempty"`
	ReportReason string `json:"reason,omitempty" db:"reason,omitempty"`
}

var db *sqlx.DB

// ADD A MUTEX?########

// Generate a random string of 32 characters
func generateRandomString() string {
	randomNum := strconv.Itoa(rand.Int())
	hash := md5.Sum([]byte(randomNum))
	return hex.EncodeToString(hash[:])
}

// Parse the URL to find the post link ID and scan the database for the corresponding entry
func getEntryForRequestedLink(w http.ResponseWriter, r *http.Request) (post, error) {
	lastIndex := strings.LastIndex(r.URL.Path, "/")
	postLinkID := r.URL.Path[lastIndex+1:]

	// Retrieve the post from the database
	query := `SELECT *
			 	FROM Posts, Links 
			 	WHERE Posts.post_id = Links.post_id and link_id = $1`
	entry := post{}
	err := db.Get(&entry, query, postLinkID)
	if err == sql.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		log.Println("error: no entries found")
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("error: unsuccessful data lookup")
	}
	return entry, err
}

// Retrieve a individual post based on the links
func handleRetrievePost(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI(), r.Method)

	entry, err := getEntryForRequestedLink(w, r)
	if err != nil {
		return
	}

	if entry.Access == "View" { // View Access
		// To avoid JSON encoding
		entry.Scope = ""
		entry.LinkID = ""

		// Encode and Send Response To Client
		err = json.NewEncoder(w).Encode(entry)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println("error: JSON encoding unsuccessful")
			return
		}

	} else { // Edit Access
		// look up the view-id
		var viewID string

		query := `SELECT link_id 
					FROM Links
					WHERE post_id = $1 and access = "View"`
		err = db.Get(&viewID, query, entry.PostID)
		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			fmt.Println("no Entries found")
			return
		} else if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("unsuccessful data lookup")
			return
		}

		lastSlashIndex := strings.LastIndex(r.URL.Path, "/")
		resourceWithoutLinkID := r.URL.Path[:lastSlashIndex]

		links := postLinks{
			EditLink: path.Join(r.Host, resourceWithoutLinkID, entry.LinkID),
			ViewLink: path.Join(r.Host, resourceWithoutLinkID, viewID),
		}
		entry.LinkID = ""

		// Combine the post and links structs for JSON encoding
		response := struct {
			post
			postLinks
		}{entry, links}

		// Encode and Send Response To Client
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print("error: JSON encoding unsuccessful")
		}
	}
}

// Retrieve all public posts
func handleRetrievePosts(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI(), r.Method)

	// Retrieve the page number from the query string
	pageNum, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println("error: invalid page number")
		return
	}

	// Retrieve the public records from the database
	limit := 5
	entries := []post{}
	offset := pageNum * limit

	query := `SELECT title, body, epoch, Links.link_id 
				FROM Posts p, Links
				WHERE Links.post_id = p.post_id and scope = "Public" and access = "View" and
				not exists
					(SELECT post_id
						FROM report
						WHERE report.post_id = p.post_id)
				ORDER BY epoch DESC
				LIMIT $1 OFFSET $2`

	err = db.Select(&entries, query, limit, offset)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("error: data retrieval unsuccessful")
		return
	}

	// Change Link ID to Links to the post
	for i := range entries {
		entries[i].LinkID = path.Join(r.Host, r.URL.Path, entries[i].LinkID)
	}

	// Encode and Send Response To Client
	err = json.NewEncoder(w).Encode(entries)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print("error: JSON encoding unsuccessful")
	} else {
		log.Println("Data Retrieval and Encoding successful")
	}
}

// Delete the post with the given link if having edit access
func handleDeletePost(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI(), r.Method)

	entry, err := getEntryForRequestedLink(w, r)
	if err != nil {
		return
	}

	if entry.Access == "View" { // Post removal forbidden
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		log.Println("error: view links cannot delete posts")
	} else { // Try post removal
		query := `PRAGMA foreign_keys = ON;
					DELETE FROM posts where post_id = $1`

		_, err := db.Exec(query, entry.PostID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println("unsuccessful database entry removal", err)
		} else {
			log.Println("Successful data removal")
		}
	}

}

// Update the post with the given link if having edit access
func handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI(), r.Method)

	entry, err := getEntryForRequestedLink(w, r)
	if err != nil {
		return
	}

	if entry.Access == "View" { // Post removal forbidden
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		log.Println("error: view links cannot update posts")
	} else { // Try post editing
		// Decode Post Contents
		updatedPostContents := post{}
		err := json.NewDecoder(r.Body).Decode(&updatedPostContents)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println("error: JSON decoding error occured")
			return
		}

		// Generate the query based on the fields passed
		query := `UPDATE Posts SET `

		if updatedPostContents.Title != "" {
			query += `title=:title, `
		}
		if updatedPostContents.Body != "" {
			query += `body=:body, `
		}
		if updatedPostContents.Scope != "" {
			query += `scope=:scope `
		}

		query += `WHERE post_id=:post_id`
		fmt.Println(query)

		updatedPostContents.PostID = entry.PostID
		_, err = db.NamedExec(query, updatedPostContents)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println("error: unsuccessful entry update", err)
		} else {
			log.Println("Data updated successfully")
		}
	}
}

// isUniqueViolation returns true if the supplied error resulted from a primary key constraint failure.
func isUniqueViolation(err error) bool {
	if err, ok := err.(sqlite3.Error); ok {
		return err.Code == 19 && err.ExtendedCode == 1555
	}
	return false
}

// Add a unique link to the database and return it
func addLinkIDToDatabase(w http.ResponseWriter, linkID string, postID int64, access string) string {
	query := `INSERT INTO links (link_id, access, post_id)
                       VALUES ($1, $2, $3)`
	_, err := db.Exec(query, linkID, access, postID)
	if err != nil {
		if isUniqueViolation(err) {
			linkID = addLinkIDToDatabase(w, generateRandomString(), postID, access)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Errorf("post creation unsuccessful")
		}
	}
	return linkID
}

// Create a post with the contents given by client and respond back with the access links
func handleCreatePost(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI(), r.Method)
	// Decode Post Contents
	newPost := post{}
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("error: decoding error occured")
		return
	}
	// Insert data in posts table
	var result sql.Result
	query := `INSERT INTO posts (title, body, scope, epoch)
                       VALUES ($1, $2, $3, $4)`
	result, err = db.Exec(query, newPost.Title, newPost.Body, newPost.Scope, time.Now().Unix())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("error: post creation unsuccessful")
		return
	}
	postID, err := result.LastInsertId()

	// Insert the Link ID's to the Link table
	editID := addLinkIDToDatabase(w, generateRandomString(), postID, "Edit")
	viewID := addLinkIDToDatabase(w, generateRandomString(), postID, "View")

	editLink := path.Join(r.Host, r.RequestURI, editID)
	viewLink := path.Join(r.Host, r.RequestURI, viewID)

	// Encode and Send Response To Client
	response := postLinks{EditLink: editLink, ViewLink: viewLink}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("error: encoding unsuccessful")
	}
}

func handlePostReport(w http.ResponseWriter, r *http.Request) {
	entry, err := getEntryForRequestedLink(w, r)
	if err != nil {
		return
	}

	if entry.Access == "Edit" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		log.Println("error: cannot report a post with edit access")
	} else {
		// Decode Post Contents
		updatedPostContents := post{}
		err := json.NewDecoder(r.Body).Decode(&updatedPostContents)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println("error: JSON decoding error occured")
			return
		}
		updatedPostContents.PostID = entry.PostID

		// Reporting the post
		query := `INSERT INTO REPORT (reason, post_id) VALUES (:reason, :post_id)`
		_, err = db.NamedExec(query, updatedPostContents)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println("error: unsuccessful post reporting ", err)
		} else {
			log.Println("Post reported successfully!")
		}
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
	r.Path("/api/v1/posts").Methods("GET").HandlerFunc(handleRetrievePosts)
	r.Path("/api/v1/posts").Methods("POST").HandlerFunc(handleCreatePost)

	// For individual post
	// ##########CHANGE IT TO PROPER REGEX FOR HEX
	r.Path("/api/v1/posts/{*}").Methods("GET").HandlerFunc(handleRetrievePost)
	r.Path("/api/v1/posts/{*}").Methods("UPDATE").HandlerFunc(handleUpdatePost)
	r.Path("/api/v1/posts/report/{*}").Methods("UPDATE").HandlerFunc(handlePostReport)
	r.Path("/api/v1/posts/{*}").Methods("DELETE").HandlerFunc(handleDeletePost)

	log.Fatal(http.ListenAndServe(":8110", r))
}
