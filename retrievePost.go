package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// Retrieve a individual post based on the links
func handleRetrievePost(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI(), r.Method)

	entry, err := getEntryForRequestedLink(w, r)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")

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
			EditLink: getRequestType(r) + path.Join(r.Host, resourceWithoutLinkID, entry.LinkID),
			ViewLink: getRequestType(r) + path.Join(r.Host, resourceWithoutLinkID, viewID),
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

// Retrieve the value of the query parameter from the query string and return the value
// and the associated error (if any)
func getQueryValue(w http.ResponseWriter, r *http.Request, queryParameter string) (int, error) {
	result, err := strconv.Atoi(r.URL.Query().Get(queryParameter))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Println("error: invalid", queryParameter)
		return 0, err
	}
	return result, nil
}

// Retrieve all public posts that are not reported
func handleRetrievePosts(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI(), r.Method)

	// Retrieve the offset value from the query string
	offset, err := getQueryValue(w, r, "offset")
	if err != nil {
		return
	}

	// Retrieve the limit value from the query string
	limit, err := getQueryValue(w, r, "limit")
	if err != nil {
		return
	}

	// Retrieve the public records from the database
	entries := []post{}

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
		entries[i].LinkID = getRequestType(r) + path.Join(r.Host, r.URL.Path, entries[i].LinkID)
	}

	w.Header().Set("Content-Type", "application/json")

	// Encode and Send Response To Client
	err = json.NewEncoder(w).Encode(entries)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print("error: JSON encoding unsuccessful")
	} else {
		log.Println("Data Retrieval and Encoding successful")
	}
}
