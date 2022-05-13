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
	log.Printf("Server now running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", logHandler{}))
}

// logHandler simply logs the incoming requests, then forwards them to
// http.DefaultServeMux.
type logHandler struct{}

// ServeHTTP implements http.Handler.
func (l logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request received: %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
	http.DefaultServeMux.ServeHTTP(w, r)
}

// HTTP HANDLER FUNCTIONS

// Handles requests to list artists in database.
func artistsHandler(w http.ResponseWriter, r *http.Request) {
	if ok := checkMethod(r.Method, []string{http.MethodGet}, w); !ok {
		return
	}

	artists, err := db.GetArtists()
	if err != nil {
		serverError(err, "could not get artists", w)
		return
	}
	writeJSON(w, artists)
}

// Handles requests to list songs by a given artist.
func songsHandler(w http.ResponseWriter, r *http.Request) {
	if ok := checkMethod(r.Method, []string{http.MethodGet}, w); !ok {
		return
	}
	if !r.URL.Query().Has("artist") {
		http.Error(w, `must provide query param "artist"`, http.StatusBadRequest)
		return
	}
	artist := r.URL.Query().Get("artist")

	songs, err := db.GetSongs(artist)
	if err != nil {
		serverError(err, "could not get songs", w)
		return
	}
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
		http.Error(w, fmt.Sprintf("invalid id %q", idstr), http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		chords, err := db.GetChords(id)
		if err != nil {
			serverError(err, "could not get chords", w)
			return
		}
		w.Write([]byte(chords))

	} else if r.Method == http.MethodPut {
		// TODO: Check for authentication
		err := db.SetChords(id, r.Body)
		if err == nil {
			// Success - nothing returned
			w.WriteHeader(http.StatusNoContent)
		} else {
			serverError(err, "could not update chords", w)
		}
	}
}

// Handles requests to create a chord sheet.
func newChordsHandler(w http.ResponseWriter, r *http.Request) {
	if ok := checkMethod(r.Method, []string{http.MethodPost}, w); !ok {
		return
	}
	// TODO: Check for authentication
	// parse request `r`
	// send to get function
	// write output to `w`
	http.Error(w, "create chords not yet implemented", http.StatusNotImplemented)
}

// Handles requests to search the database for a song.
func searchHandler(w http.ResponseWriter, r *http.Request) {
	if ok := checkMethod(r.Method, []string{http.MethodGet}, w); !ok {
		return
	}
	// parse request `r`
	// send to update function
	// write output to `w`
	http.Error(w, "search not yet implemented", http.StatusNotImplemented)
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
		), http.StatusMethodNotAllowed)
	}
	return allow
}

// serverError returns a 500 response, and logs the offending error.
func serverError(e error, msg string, w http.ResponseWriter) {
	log.Printf("ERROR: %v", e)
	http.Error(w, msg, http.StatusInternalServerError)
}

// writeJSON marshals `data` to JSON and writes it to `w`.
func writeJSON(w http.ResponseWriter, data interface{}) {
	jData, err := json.Marshal(data)
	if err != nil {
		serverError(err, "error marshalling to JSON", w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}
