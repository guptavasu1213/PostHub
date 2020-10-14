package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
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
