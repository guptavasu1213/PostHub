// Vasu Gupta
// ID: 3066521
// Assignment 2

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
)

// Serves the page containing all the posts through the response writer
func handlerToRetrieveAllPostsPage(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI(), r.Method)

	// Creating HTML template
	template, err := template.ParseFiles("dist/templates/allPosts.tmpl", navBarTemplatePath, footerTemplatePath)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": Error in template parsing", http.StatusInternalServerError)
		log.Println("error: could not generate html template")
		return
	}

	// Add struct values to the template
	template.Execute(w, nil)
}

// Extract view ID for a post using Post ID
func getViewIDFromPostID(w http.ResponseWriter, postID int64) (string, error) {
	var viewID string

	query := `SELECT link_id 
						FROM Links
						WHERE post_id = $1 and access = "View"`
	err := db.Get(&viewID, query, postID)
	if err == sql.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		fmt.Println("no Entries found")
		return "", err
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("unsuccessful data lookup")
		return "", err
	}
	return viewID, nil
}

// Serve HTML template of the View Only page to the client
func retrieveReadonlyPostPage(w http.ResponseWriter, entry post) {
	// Creating HTML template
	template, err := template.ParseFiles("dist/templates/publicPortal.tmpl", navBarTemplatePath, footerTemplatePath)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": Error in template parsing", http.StatusInternalServerError)
		log.Println("error: could not generate html template")
		return
	}

	// Add struct values to the template
	template.Execute(w, entry)
}

// Serve HTML template of the Admin Portal to the client
func retrieveAdminPostPage(w http.ResponseWriter, r *http.Request, entry post) {
	// Creating HTML template
	template, err := template.ParseFiles("dist/templates/adminPortal.tmpl", navBarTemplatePath, footerTemplatePath)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": Error in template parsing", http.StatusInternalServerError)
		log.Println("error: could not generate html template")
		return
	}

	// Retrieving the view id of the post
	viewID, err := getViewIDFromPostID(w, entry.PostID)
	if err != nil {
		return
	}

	links := postLinks{
		EditLink: path.Join(r.Host, "posts", entry.LinkID),
		ViewLink: path.Join(r.Host, "posts", viewID),
	}

	// Combine the post and links structs
	combinedStruct := struct {
		post
		postLinks
	}{entry, links}

	// Add struct values to the template
	template.Execute(w, combinedStruct)
}

// Handler for serving HTML template of the post to the client
func handleIndividualPageServing(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.RequestURI(), r.Method)

	entry, err := getEntryForRequestedLink(w, r)
	if err != nil {
		return
	}

	if entry.Access == "View" {
		log.Println("View Access")
		retrieveReadonlyPostPage(w, entry)

	} else { // Edit Access
		log.Println("Edit Access")
		retrieveAdminPostPage(w, r, entry)
	}
}

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
		viewID, err := getViewIDFromPostID(w, entry.PostID)
		if err != nil {
			return
		}

		links := postLinks{
			EditLink: entry.LinkID,
			ViewLink: viewID,
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
