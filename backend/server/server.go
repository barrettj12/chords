// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// backend/server/server.go
// Contains the app's HTTP server and request handlers

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/barrettj12/chords/backend/dblayer"
)

type Server struct {
	db   dblayer.ChordsDB
	addr string
}

func New(db dblayer.ChordsDB, addr string) Server {
	return Server{db, addr}
}

func (s *Server) Start() {
	// Register API endpoints
	http.HandleFunc("/artists", s.artistsHandler)  // list artists in database
	http.HandleFunc("/songs", s.songsHandler)      // list songs by a given artist
	http.HandleFunc("/chords/", s.chordsHandler)   // view/update a chord sheet
	http.HandleFunc("/chords", s.newChordsHandler) // create a chord sheet
	http.HandleFunc("/search", s.searchHandler)    // search database for song

	// Start listening on port 8080
	log.Printf(fmt.Sprintf("Server now running at http://localhost%s", s.addr))
	log.Fatal(http.ListenAndServe(s.addr, handler{}))
}

// logHandler does some extra post-request / pre-response handling common
// to all requests - see the ServeHTTP method below.
type handler struct{}

// ServeHTTP implements http.Handler.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log request
	log.Printf("request: %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
	// Add CORS header
	w.Header().Set("Access-Control-Allow-Origin", "*")

	http.DefaultServeMux.ServeHTTP(w, r)
}

// HTTP HANDLER FUNCTIONS

// Handles requests to list artists in database.
func (s *Server) artistsHandler(w http.ResponseWriter, r *http.Request) {
	if ok := checkMethod(r.Method, []string{http.MethodGet}, w); !ok {
		return
	}

	artists, err := s.db.GetArtists()
	if err != nil {
		serverError(err, "could not get artists", w)
		return
	}
	writeJSON(w, artists)
}

// Handles requests to list songs by a given artist.
func (s *Server) songsHandler(w http.ResponseWriter, r *http.Request) {
	if ok := checkMethod(r.Method, []string{http.MethodGet}, w); !ok {
		return
	}
	if !r.URL.Query().Has("artist") {
		http.Error(w, `must provide query param "artist"`, http.StatusBadRequest)
		return
	}
	artist := r.URL.Query().Get("artist")

	songs, err := s.db.GetSongs(artist)
	if err != nil {
		serverError(err, "could not get songs", w)
		return
	}
	writeJSON(w, songs)
}

// Handles requests to view/update a chord sheet.
func (s *Server) chordsHandler(w http.ResponseWriter, r *http.Request) {
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
		chords, err := s.db.GetChords(id)
		if err != nil {
			serverError(err, "could not get chords", w)
			return
		}
		w.Write([]byte(chords))

	} else if r.Method == http.MethodPut {
		// TODO: Check for authentication
		chords, err := io.ReadAll(r.Body)
		if err != nil {
			serverError(err, "could not update chords", w)
			return
		}
		err = s.db.SetChords(id, chords)
		if err != nil {
			serverError(err, "could not update chords", w)
			return
		}
		// Success - nothing returned
		w.WriteHeader(http.StatusNoContent)
	}
}

// Handles requests to create a chord sheet.
func (s *Server) newChordsHandler(w http.ResponseWriter, r *http.Request) {
	if ok := checkMethod(r.Method, []string{http.MethodPost}, w); !ok {
		return
	}
	// TODO: Check for authentication

	body, err := io.ReadAll(r.Body)
	if err != nil {
		serverError(err, "could not update chords", w)
		return
	}
	// Unmarshal r.Body from JSON
	nc := &dblayer.NewChords{}
	json.Unmarshal(body, nc)

	id, err := s.db.MakeChords(*nc)
	if err != nil {
		serverError(err, "could not create chords", w)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Location", fmt.Sprintf("/chords/%d", id))
}

// Handles requests to search the database for a song.
func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
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
