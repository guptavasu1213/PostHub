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
	"strconv"
	"time"

	sqlite3 "github.com/mattn/go-sqlite3"
)

// Generate a random string of 32 characters
func generateRandomString() string {
	randomNum := strconv.Itoa(rand.Int())
	hash := md5.Sum([]byte(randomNum))
	return hex.EncodeToString(hash[:])
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
		err := json.NewDecoder(r.Body).Decode(&entry)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			log.Println("error: JSON decoding error occured")
			return
		}

		// Error Check Null Fields
		if entry.Title == "" || entry.Body == "" || entry.Scope == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			log.Println("error: all fields are required and should be non-null")
			return
		}

		// Generate the query based on the fields passed
		query := `UPDATE Posts 
					SET title=:title, body=:body, scope=:scope
					WHERE post_id=:post_id`

		_, err = db.NamedExec(query, entry)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Println("error: unsuccessful entry update")
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
			fmt.Println("post creation unsuccessful")
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
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println("error: decoding error occured", err)
		return
	}
	// Error Check Null Fields
	if newPost.Title == "" || newPost.Body == "" || newPost.Scope == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println("error: all fields are required and should be non-null")
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

	editLink := editID
	viewLink := viewID

	// Encode and Send Response To Client
	response := postLinks{EditLink: editLink, ViewLink: viewLink}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("error: encoding unsuccessful")
	}
}

// Report the post based on the link
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
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			log.Println("error: JSON decoding error occured")
			return
		}
		updatedPostContents.PostID = entry.PostID

		// Error Check Null Field
		if updatedPostContents.ReportReason == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			log.Println("error: all fields are required and should be non-null")
			return
		}

		// Report the post
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
