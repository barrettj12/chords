package tests

import (
	"log"
	"os"
	"testing"

	"github.com/barrettj12/chords/src/client"
	"github.com/barrettj12/chords/src/dblayer"
	"github.com/barrettj12/chords/src/server"
	"github.com/stretchr/testify/assert"
)

func TestNewSong(t *testing.T) {
	_, _, c, teardown := setup(t)
	defer teardown()

	// Add new song via API
	newSong := dblayer.SongMeta{
		ID:       "BananaPancakes",
		Name:     "Banana Pancakes",
		Artist:   "Jack Johnson",
		Album:    "In Between Dreams",
		TrackNum: 3,
	}
	resp, err := c.NewSong(newSong)
	handleClientError(t, err)
	assert.Equal(t, resp, newSong)

	// TODO: check db state via fs?
}

func TestUpdateChords(t *testing.T) {
	db, _, c, teardown := setup(t)
	defer teardown()

	// Add new song to DB
	newSong := dblayer.SongMeta{
		ID:       "BananaPancakes",
		Name:     "Banana Pancakes",
		Artist:   "Jack Johnson",
		Album:    "In Between Dreams",
		TrackNum: 3,
	}
	retSong, err := db.NewSong(newSong)
	assert.Nil(t, err)
	assert.Equal(t, retSong, newSong)

	// Update chords via API
	chords := dblayer.Chords(`
intro/chorus:
Am7 - G7

verse:
(D7) - G7 - D7 - Am7 - C7

bridge:
Am7 - D     (x2)
Bm7 - Em7 - (D+maj7) - C
G - D7 - G
`)
	resp, err := c.UpdateChords(newSong.ID, chords)
	handleClientError(t, err)
	assert.EqualValues(t, resp, chords)

	// Get chords from API
	retChords, err := c.GetChords(newSong.ID)
	handleClientError(t, err)
	assert.EqualValues(t, retChords, chords)

	// TODO: check db state via fs?
}

func setup(t *testing.T) (dblayer.ChordsDB, *server.Server, *client.Client, func()) {
	// Set up DB
	dataDir, err := os.MkdirTemp("", "data")
	assert.Nil(t, err)
	logger := log.Default()
	db := dblayer.NewLocalfs(dataDir, logger)

	// Set up server
	authKey := "passwordfoo"
	s, err := server.New(db, ":8080", logger, authKey)
	assert.Nil(t, err)
	go func() {
		err = s.Run()
		assert.Nil(t, err)
	}()

	// Set up client
	c, err := client.NewClient("http://localhost:8080", authKey)
	assert.Nil(t, err)

	teardown := func() {
		// Remove tempdir for DB
		err := os.RemoveAll(dataDir)
		assert.Nil(t, err)

		// Kill server
		err = s.Kill()
		assert.Nil(t, err)
	}

	return db, s, c, teardown
}

func handleClientError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}
}
