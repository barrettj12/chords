package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/barrettj12/chords/src/dblayer"
	"github.com/stretchr/testify/assert"
)

func TestArtists(t *testing.T) {
	// Set up DB & server, http writer
	db := dblayer.NewTempDB()
	s := Server{api: &ChordsAPI{db: db}}

	// Add artists to DB
	artists := []string{"Elton John", "Rod Stewart", "Spacehog", "foobar"}
	for _, a := range artists {
		db.NewSong(dblayer.SongMeta{
			Artist: a,
		})
	}

	// Get artists via API
	r := httptest.NewRequest(http.MethodGet, "/api/v0/songs", nil)
	w := httptest.NewRecorder()
	s.api.artistsHandler(w, r)

	res := w.Result()
	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK, "body: %s", data)

	respArtists := []string{}
	err = json.Unmarshal(data, &respArtists)
	assert.Nil(t, err)
	assert.ElementsMatch(t, respArtists, artists)
}

func TestNewSong(t *testing.T) {
	// Set up DB
	dataDir, err := os.MkdirTemp("", "data")
	assert.Nil(t, err)
	defer func() {
		err := os.RemoveAll(dataDir)
		assert.Nil(t, err)
	}()

	logger := log.Default()
	db := dblayer.NewLocalfs(dataDir, logger)

	s := Server{api: &ChordsAPI{db: db}}
	w := httptest.NewRecorder()

	// Add new song via API
	newSong := dblayer.SongMeta{
		ID:       "BananaPancakes",
		Name:     "Banana Pancakes",
		Artist:   "Jack Johnson",
		Album:    "In Between Dreams",
		TrackNum: 3,
	}

	data, err := json.Marshal(newSong)
	assert.Nil(t, err)
	body := bytes.NewReader(data)
	r := httptest.NewRequest(http.MethodPost, "/api/v0/songs", body)

	// API call
	s.api.newSong(w, r)
	res := w.Result()
	data, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK, "body: %s", data)

	// Check response body matches input data
	respMeta := dblayer.SongMeta{}
	err = json.Unmarshal(data, &respMeta)
	assert.Nil(t, err)
	assert.Equal(t, newSong, respMeta)

	// TODO: check db state via fs?
}

func TestUpdateSong(t *testing.T) {
	// Set up DB, server, http writer
	db := dblayer.NewTempDB()
	s := Server{api: &ChordsAPI{db: db}}
	w := httptest.NewRecorder()

	// Put a song in the database
	initMeta := dblayer.SongMeta{
		ID:       "BananaPancakes",
		Name:     "Banana Panckaes",
		Artist:   "Jack Jonhson",
		Album:    "In Bewteen Dreams",
		TrackNum: 4,
	}
	newMeta, err := db.NewSong(initMeta)
	assert.Nil(t, err)
	id := newMeta.ID

	// Update metadata via API
	updatedMeta := dblayer.SongMeta{
		ID:       id,
		Name:     "Banana Pancakes",
		Artist:   "Jack Johnson",
		Album:    "In Between Dreams",
		TrackNum: 3,
	}
	data, err := json.Marshal(updatedMeta)
	assert.Nil(t, err)
	body := bytes.NewReader(data)
	r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v0/songs?id=%s", id), body)

	// API call
	s.api.updateSong(w, r)
	res := w.Result()
	data, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK, "body: %s", data)

	// Check response body matches input data
	respMeta := dblayer.SongMeta{}
	err = json.Unmarshal(data, &respMeta)
	assert.Nil(t, err)
	assert.Equal(t, updatedMeta, respMeta)

	// Check db state
	dbSongs, err := db.GetSongs("", id)
	assert.Nil(t, err)
	assert.Len(t, dbSongs, 1)
	assert.Equal(t, dbSongs[0], respMeta)
}
