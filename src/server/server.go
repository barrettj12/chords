// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// src/server/server.go
// Contains the app's HTTP server and request handlers

package server

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"

	gqlhandle "github.com/99designs/gqlgen/graphql/handler"
	gqlplay "github.com/99designs/gqlgen/graphql/playground"
	"github.com/barrettj12/chords/gqlgen"
	"github.com/barrettj12/chords/src/data"
	"github.com/barrettj12/chords/src/dblayer"
)

type Server struct {
	httpServer http.Server
	listener   net.Listener
	logger     *log.Logger
	api        *ChordsAPI
}

// New returns a new Server with the specified DB and address. `logFlags` is
// as provided to log.New - see https://pkg.go.dev/log#pkg-constants
func New(db dblayer.ChordsDB, addr string, logger *log.Logger, authKey string) (*Server, error) {
	frontend, err := NewFrontend(fmt.Sprintf("http://localhost%s", addr))
	if err != nil {
		return nil, err
	}
	api := ChordsAPI{db, logger, authKey}

	return &Server{
		httpServer: http.Server{
			Addr: addr,
			Handler: newHandler(
				logger,
				&api,
				frontend,
			),
		},
		logger: logger,
		api:    &api,
	}, nil
}

// Listen opens a network connection (non-blocking) and returns the address
// that it's listening on.
func (s *Server) Listen() (net.Addr, error) {
	// copied from net/http.Server.ListenAndServe
	addr := s.httpServer.Addr
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s.logger.Printf("server now listening at %s", ln.Addr())
	s.listener = ln
	return ln.Addr(), nil
}

// Serve serves the HTTP server (blocking)
func (s *Server) Serve() error {
	closeErr := s.httpServer.Serve(s.listener)
	if errors.Is(closeErr, http.ErrServerClosed) {
		s.logger.Println("server closed")
		return nil
	}
	return closeErr
}

func (s *Server) Run() error {
	_, err := s.Listen()
	if err != nil {
		return err
	}
	return s.Serve()
}

// Shuts down the HTTP server - necessary for running tests back-to-back.
func (s *Server) Kill() error {
	return s.httpServer.Shutdown(context.Background())
}

// handler does some extra post-request / pre-response handling common
// to all requests - see the ServeHTTP method below.
type handler struct {
	logger *log.Logger
	mux    *http.ServeMux
}

func newHandler(logger *log.Logger, api *ChordsAPI, frontend *Frontend) handler {
	// Set up mux
	mux := http.NewServeMux()

	// Register API endpoints
	mux.HandleFunc("/api/v0/artists", api.artistsHandler)  // list artists in database
	mux.HandleFunc("/api/v0/songs", api.songsHandler)      // song metadata API
	mux.HandleFunc("/api/v0/chords", api.chordsHandler)    // view/update a chord sheet
	mux.HandleFunc("/api/v0/see-also", api.seeAlsoHandler) // get related artists
	mux.HandleFunc("/api/v0/random", api.randomHandler)    // get random chords
	mux.HandleFunc("/api/v0/search", api.searchHandler)    // search chords

	// Favicon
	mux.HandleFunc("/favicon.ico", serveFavicon)

	// Register frontend endpoints
	frontend.registerHandlers(mux)

	// GraphQL endpoints
	mux.Handle("/graphql", gqlhandle.NewDefaultServer(gqlgen.NewExecutableSchema(gqlgen.Config{
		Resolvers: &gqlgen.Resolver{
			DB: data.AsChordsDBv1(api.db),
		}})))
	mux.Handle("/graphql/playground", gqlplay.Handler("GraphQL playground", "/graphql"))

	return handler{
		logger: logger,
		mux:    mux,
	}
}

// ServeHTTP implements http.Handler.
func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log request
	h.logger.Printf("request: %s %s from %s\n", r.Method, r.RequestURI, r.RemoteAddr)

	// Log body (for debugging)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("ERROR reading request body: %s\n", err)
		http.Error(w, "reading request body", http.StatusInternalServerError)
		return
	}

	r.Body.Close()
	h.logger.Printf("request body:\n%s\n", body)
	r.Body = io.NopCloser(bytes.NewReader(body))

	// Duplicate writer so we can view response
	// Only want to log API responses
	if strings.HasPrefix(r.URL.Path, "/api") {
		w2 := NewResponseWriterWrapper(w)
		defer func() {
			resp := w2.String()
			h.logger.Printf("sending response:\n%s\n", resp)
		}()
		w = w2
	}

	// Add CORS header
	w.Header().Set("Access-Control-Allow-Origin", "*")
	h.mux.ServeHTTP(w, r)
}

// API HANDLERS

type ChordsAPI struct {
	db      dblayer.ChordsDB
	logger  *log.Logger
	authKey string
}

// Handles requests to the /api/v0/artists endpoint.
func (s *ChordsAPI) artistsHandler(w http.ResponseWriter, r *http.Request) {
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
func (s *ChordsAPI) songsHandler(w http.ResponseWriter, r *http.Request) {
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
func (s *ChordsAPI) getSongs(w http.ResponseWriter, r *http.Request) {
	artist := r.URL.Query().Get("artist")
	id := r.URL.Query().Get("id")
	query := r.URL.Query().Get("query")
	songs, err := s.db.GetSongs(artist, id, query)

	if err == nil {
		s.writeJSON(w, songs)
	} else {
		s.serverError(err, "could not get songs", w)
	}
}

// Add a new song to the database.
func (s *ChordsAPI) newSong(w http.ResponseWriter, r *http.Request) {
	if !s.authorised(w, r) {
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
func (s *ChordsAPI) updateSong(w http.ResponseWriter, r *http.Request) {
	if !s.authorised(w, r) {
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
func (s *ChordsAPI) deleteSong(w http.ResponseWriter, r *http.Request) {
	if !s.authorised(w, r) {
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
func (s *ChordsAPI) chordsHandler(w http.ResponseWriter, r *http.Request) {
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
func (s *ChordsAPI) getChords(w http.ResponseWriter, r *http.Request) {
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
func (s *ChordsAPI) updateChords(w http.ResponseWriter, r *http.Request) {
	if !s.authorised(w, r) {
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

// Handles requests to the /api/v0/see-also endpoint.
func (s *ChordsAPI) seeAlsoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// TODO: allow updating see also data via POST/PUT/PATCH/DELETE
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	artist := r.URL.Query().Get("artist")
	relatedArtists, err := s.db.SeeAlso(artist)
	if err != nil {
		s.serverError(err, "could not get related artists", w)
		return
	}

	s.writeJSON(w, relatedArtists)
}

func (s *ChordsAPI) randomHandler(w http.ResponseWriter, r *http.Request) {
	allSongs, err := s.db.GetSongs("", "", "")
	if err != nil {
		s.serverError(err, "getting songs", w)
		return
	}

	n := rand.Intn(len(allSongs))
	s.writeJSON(w, allSongs[n])
}

func (s *ChordsAPI) searchHandler(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("q") {
		http.Error(w, `missing query param "q"`, http.StatusBadRequest)
		return
	}

	searchQuery := r.URL.Query().Get("q")
	if searchQuery == "" {
		http.Error(w, `search query cannot be empty`, http.StatusBadRequest)
		return
	}

	results, err := s.db.Search(searchQuery)
	if err != nil {
		s.serverError(err, "getting songs", w)
		return
	}

	s.writeJSON(w, results)
}

//go:embed favicon.ico
var faviconData []byte

// Serve favicon data in favicon.go.
func serveFavicon(w http.ResponseWriter, _ *http.Request) {
	w.Write(faviconData)
}

// HELPER FUNCTIONS

// For methods which write to the database, check we are authorised to do this.
// If not, write an Unauthorised error to w.
func (s *ChordsAPI) authorised(w http.ResponseWriter, r *http.Request) bool {
	// TODO: work out how to authorise
	key := r.Header.Get("Authorization")
	authd := key == s.authKey

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
func (s *ChordsAPI) serverError(e error, msg string, w http.ResponseWriter) {
	s.logger.Printf("ERROR: %v", e)
	http.Error(w, msg, http.StatusInternalServerError)
}

// writeJSON marshals `data` to JSON and writes it to `w`.
func (s *ChordsAPI) writeJSON(w http.ResponseWriter, data any) {
	jData, err := json.Marshal(data)
	if err != nil {
		s.serverError(err, "error marshalling to JSON", w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

// ResponseWriter wrapper allowing inspection of the outgoing response
// Credit to Alessandro Argentieri on Stack Overflow
// https://stackoverflow.com/a/66531582

// ResponseWriterWrapper struct is used to log the response
type ResponseWriterWrapper struct {
	w          *http.ResponseWriter
	body       *bytes.Buffer
	statusCode *int
}

// NewResponseWriterWrapper static function creates a wrapper for the http.ResponseWriter
func NewResponseWriterWrapper(w http.ResponseWriter) ResponseWriterWrapper {
	var buf bytes.Buffer
	var statusCode int = 200
	return ResponseWriterWrapper{
		w:          &w,
		body:       &buf,
		statusCode: &statusCode,
	}
}

func (rww ResponseWriterWrapper) Write(buf []byte) (int, error) {
	rww.body.Write(buf)
	return (*rww.w).Write(buf)
}

// Header function overwrites the http.ResponseWriter Header() function
func (rww ResponseWriterWrapper) Header() http.Header {
	return (*rww.w).Header()
}

// WriteHeader function overwrites the http.ResponseWriter WriteHeader() function
func (rww ResponseWriterWrapper) WriteHeader(statusCode int) {
	(*rww.statusCode) = statusCode
	(*rww.w).WriteHeader(statusCode)
}

func (rww ResponseWriterWrapper) String() string {
	return rww.body.String()
}
