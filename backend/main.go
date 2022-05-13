// add license info here

// add description of file here

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/barrettj12/jordys-chordies/backend/dblayer"
)

// Database
var db dblayer.ChordsDB

func main() {
	// Register API endpoints
	http.HandleFunc("/artists", artistsHandler)  // list artists in database
	http.HandleFunc("/songs", songsHandler)      // list songs by a given artist
	http.HandleFunc("/chords/", chordsHandler)   // view/update a chord sheet
	http.HandleFunc("/chords", newChordsHandler) // create a chord sheet
	http.HandleFunc("/search", searchHandler)    // search database for song

	// Start listening on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// HTTP HANDLERS

// Handles requests to list artists in database.
func artistsHandler(w http.ResponseWriter, r *http.Request) {
	if ok := checkMethod(r.Method, []string{http.MethodGet}, w); !ok {
		return
	}
	artists := db.GetArtists()
	writeJSON(w, artists)
}

// Handles requests to list songs by a given artist.
func songsHandler(w http.ResponseWriter, r *http.Request) {
	if ok := checkMethod(r.Method, []string{http.MethodGet}, w); !ok {
		return
	}
	if !r.URL.Query().Has("artist") {
		http.Error(w, `must provide query param "artist"`, 400)
		return
	}
	artist := r.URL.Query().Get("artist")
	songs := db.GetSongs(artist)
	writeJSON(w, songs)
}

// Handles requests to view/update a chord sheet.
func chordsHandler(w http.ResponseWriter, r *http.Request) {
	if ok := checkMethod(r.Method, []string{http.MethodGet, http.MethodPut}, w); !ok {
		return
	}
	idstr := r.URL.Path[8:] // 8 = len("/chords/")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid id %q", idstr), 400)
		return
	}

	if r.Method == http.MethodGet {
		chords := db.GetChords(id)
		w.Write([]byte(chords))
	} else if r.Method == http.MethodPut {
		// Check for authentication
		// should put return the updated chords?
	}
}

// Handles requests to create a chord sheet.
func newChordsHandler(w http.ResponseWriter, r *http.Request) {
	if ok := checkMethod(r.Method, []string{http.MethodPost}, w); !ok {
		return
	}
	// Check for authentication
	// parse request `r`
	// send to get function
	// write output to `w`
}

// Handles requests to search the database for a song.
func searchHandler(w http.ResponseWriter, r *http.Request) {
	if ok := checkMethod(r.Method, []string{http.MethodGet}, w); !ok {
		return
	}
	// parse request `r`
	// send to update function
	// write output to `w`
	w.Write([]byte("search not yet implemented"))
}

// HELPER FUNCTIONS

// Returns whether the attempted method `method` is allowed by this endpoint
// (i.e. in the slice `allowed`). If not, it will also write
// "405 Method Not Allowed" to `w`.
func checkMethod(method string, allowed []string, w http.ResponseWriter) bool {
	allow := false
	for _, m := range allowed {
		if method == m {
			allow = true
		}
	}

	if !allow && w != nil {
		http.Error(w, fmt.Sprintf(
			"method %s not allowed; allowed methods are %s",
			method, strings.Join(allowed, ", "),
		), 405)
	}
	return allow
}

// serverError returns a 500 response.
func serverError(w http.ResponseWriter) {
	http.Error(w, "", 500)
}

// writeJSON marshals `data` to JSON and writes it to `w`.
func writeJSON(w http.ResponseWriter, data interface{}) {
	jData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "error marshalling to JSON", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}
