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
	"os"

	"github.com/barrettj12/chords/backend/dblayer"
)

type Server struct {
	db   dblayer.ChordsDB
	addr string
	log  *log.Logger
}

// New returns a new Server with the specified DB and address. `logFlags` is
// as provided to log.New - see https://pkg.go.dev/log#pkg-constants
func New(db dblayer.ChordsDB, addr string, logFlags int) Server {
	return Server{db, addr, log.New(os.Stdout, "", logFlags)}
}

func (s *Server) Start() {
	// Register API endpoints
	http.HandleFunc("/api/v0/artists", s.artistsHandler) // list artists in database
	http.HandleFunc("/api/v0/songs", s.songsHandler)     // song metadata API
	http.HandleFunc("/api/v0/chords", s.chordsHandler)   // view/update a chord sheet

	// Start listening on port 8080
	s.log.Printf(fmt.Sprintf("Server now running at http://localhost%s", s.addr))
	s.log.Fatal(http.ListenAndServe(s.addr, handler{s.log}))
}

// handler does some extra post-request / pre-response handling common
// to all requests - see the ServeHTTP method below.
type handler struct {
	log *log.Logger
}

// ServeHTTP implements http.Handler.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log request
	h.log.Printf("request: %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
	// Add CORS header
	w.Header().Set("Access-Control-Allow-Origin", "*")

	http.DefaultServeMux.ServeHTTP(w, r)
}

// HTTP HANDLER FUNCTIONS

// Handles requests to the /api/v0/artists endpoint.
func (s *Server) artistsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		artists, err := s.db.GetArtists()
		if err == nil {
			s.writeJSON(w, artists)
		} else {
			s.serverError(err, "could not get artists", w)
		}

	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// Handles requests to the /api/v0/songs endpoint.
func (s *Server) songsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getSongs(w, r)
	case http.MethodPost:
		s.newSong(w, r)
	case http.MethodPut:
		s.updateSong(w, r)
	case http.MethodDelete:
		s.deleteSong(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// Get metadata for songs matching query.
func (s *Server) getSongs(w http.ResponseWriter, r *http.Request) {
	artist := r.URL.Query().Get("artist")
	id := r.URL.Query().Get("id")
	// TODO: support search string parameter
	songs, err := s.db.GetSongs(artist, id)

	if err == nil {
		s.writeJSON(w, songs)
	} else {
		s.serverError(err, "could not get songs", w)
	}
}

// Add a new song to the database.
func (s *Server) newSong(w http.ResponseWriter, r *http.Request) {
	if !authorised(w, r) {
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		s.serverError(err, "io error", w)
	}

	song := &dblayer.SongMeta{}
	err = json.Unmarshal(data, song)
	if err != nil {
		s.serverError(err, "parsing body", w)
	}

	newSong, err := s.db.NewSong(*song)
	if err == nil {
		s.writeJSON(w, newSong)
	} else {
		s.serverError(err, "creating new song", w)
	}
}

// Update the metadata for a song in the database.
func (s *Server) updateSong(w http.ResponseWriter, r *http.Request) {
	if !authorised(w, r) {
		return
	}
	id, ok := idParam(w, r)
	if !ok {
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		s.serverError(err, "io error", w)
	}

	meta := &dblayer.SongMeta{}
	err = json.Unmarshal(data, meta)
	if err != nil {
		s.serverError(err, "parsing body", w)
	}

	newMeta, err := s.db.UpdateSong(id, *meta)
	if err == nil {
		s.writeJSON(w, newMeta)
	} else {
		s.serverError(err, "updating song metadata", w)
	}
}

// Delete a song from the database. The song's chords will also be deleted.
func (s *Server) deleteSong(w http.ResponseWriter, r *http.Request) {
	if !authorised(w, r) {
		return
	}
	id, ok := idParam(w, r)
	if !ok {
		return
	}

	err := s.db.DeleteSong(id)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		s.serverError(err, "deleting song", w)
	}
}

// Handles requests to the /api/v0/chords endpoint.
func (s *Server) chordsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getChords(w, r)
	case http.MethodPut:
		s.updateChords(w, r)
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// Get chords for a given song.
func (s *Server) getChords(w http.ResponseWriter, r *http.Request) {
	id, ok := idParam(w, r)
	if !ok {
		return
	}

	chords, err := s.db.GetChords(id)
	if err == nil {
		w.Write(chords)
	} else {
		s.serverError(err, "getting chords", w)
	}
}

// Update chords for a given song.
func (s *Server) updateChords(w http.ResponseWriter, r *http.Request) {
	if !authorised(w, r) {
		return
	}
	id, ok := idParam(w, r)
	if !ok {
		return
	}

	chords, err := io.ReadAll(r.Body)
	if err != nil {
		s.serverError(err, "io error", w)
	}

	newChords, err := s.db.UpdateChords(id, chords)
	if err == nil {
		w.Write(newChords)
	} else {
		s.serverError(err, "updating chords", w)
	}
}

// HELPER FUNCTIONS

// For methods which write to the database, check we are authorised to do this.
// If not, write an Unauthorised error to w.
func authorised(w http.ResponseWriter, r *http.Request) bool {
	// TODO: work out how to authorise
	authd := true

	if !authd {
		http.Error(w, "", http.StatusUnauthorized)
	}
	return authd
}

// Check if the required ID param has been provided.
// If not, write out an error
// Return id and whether it was defined.
func idParam(w http.ResponseWriter, r *http.Request) (string, bool) {
	if !r.URL.Query().Has("id") {
		http.Error(w, `required param "id" not provided`, http.StatusBadRequest)
		return "", false
	}
	return r.URL.Query().Get("id"), true
}

// serverError returns a 500 response, and logs the offending error.
func (s *Server) serverError(e error, msg string, w http.ResponseWriter) {
	s.log.Printf("ERROR: %v", e)
	http.Error(w, msg, http.StatusInternalServerError)
}

// writeJSON marshals `data` to JSON and writes it to `w`.
func (s *Server) writeJSON(w http.ResponseWriter, data any) {
	jData, err := json.Marshal(data)
	if err != nil {
		s.serverError(err, "error marshalling to JSON", w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}
