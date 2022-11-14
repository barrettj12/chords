package tests

import (
	"log"
	"os"
	"testing"

	"github.com/barrettj12/chords/backend/dblayer"
	"github.com/barrettj12/chords/backend/server"
	"github.com/barrettj12/chords/client"
	"github.com/stretchr/testify/assert"
)

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

	// Set up server
	authKey := "passwordfoo"
	s, err := server.New(db, ":8080", logger, authKey)
	assert.Nil(t, err)
	go s.Start()

	// Set up client
	c, err := client.NewClient("http://localhost:8080", authKey)
	assert.Nil(t, err)

	// Add new song via API
	newSong := dblayer.SongMeta{
		ID:       "BananaPancakes",
		Name:     "Banana Pancakes",
		Artist:   "Jack Johnson",
		Album:    "In Between Dreams",
		TrackNum: 3,
	}
	resp, err := c.NewSong(newSong)
	assert.Nil(t, err)
	assert.Equal(t, resp, newSong)

	// TODO: check db state via fs?
}
