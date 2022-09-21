// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// backend/dblayer/interface.go
// This package contains the database layer, which is responsible for
// communicating with the database, and converting between the database's
// format and Go objects.
// This file contains interfaces/code shared between different DBs

package dblayer

import (
	"fmt"
	"math/rand"
)

type ChordsDB interface {
	GetArtists() ([]string, error)
	GetSongs(artist, id string) ([]SongMeta, error)
	NewSong(SongMeta) (SongMeta, error)
	UpdateSong(id string, meta SongMeta) (SongMeta, error)
	DeleteSong(id string) error
	GetChords(id string) (Chords, error)
	UpdateChords(id string, chords Chords) (Chords, error)
	// Close() error
}

type SongMeta struct {
	id       string `json:"id"`
	name     string `json:"name"`
	artist   string `json:"artist"`
	album    string `json:"album"`
	trackNum int    `json:"id"`
}

type Chords []byte

// Fill fills a ChordsDB with some sample data - good for demonstration
// and/or testing.
func Fill(db ChordsDB) error {
	for ltr := 'A'; ltr <= 'Z'; ltr++ {
		numArtists := rand.Intn(4)
		for art := 0; art <= numArtists; art++ {
			numAlbums := rand.Intn(4)
			for alb := 0; alb <= numAlbums; alb++ {
				numSongs := rand.Intn(10)
				for sng := 0; sng <= numSongs; sng++ {
					_, err := db.NewSong(SongMeta{
						name:     fmt.Sprintf("song%d", sng),
						artist:   fmt.Sprintf("%cartist%d", ltr, art),
						album:    fmt.Sprintf("album%d", alb),
						trackNum: sng + 1,
					})
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
