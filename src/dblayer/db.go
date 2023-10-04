// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// src/dblayer/db.go
// This file contains the ChordsDB interface, which is implemented by the
// different DB "providers" in this package. Implementations are responsible
// for communicating with the database, and converting between the database's
// format and Go objects.
// It also contains the core data structures for the app (SongMeta, Chords).

package dblayer

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/barrettj12/chords/src/types"
)

type ChordsDB interface {
	GetArtists() ([]string, error)
	GetSongs(artist, id, query string) ([]SongMeta, error)
	NewSong(SongMeta) (SongMeta, error)
	UpdateSong(id string, meta SongMeta) (SongMeta, error)
	DeleteSong(id string) error
	GetChords(id string) (Chords, error)
	UpdateChords(id string, chords Chords) (Chords, error)
	SeeAlso(artist string) ([]string, error)
	Search(query string) ([]SongMeta, error)
	// Close() error
}

func GetDB(url string, logger *log.Logger) (ChordsDB, error) {
	if strings.HasPrefix(url, "postgres") {
		logger.Printf("Using Postgres database at %s\n", url)
		return NewPostgres(url)
	} else if url == "" {
		logger.Println("Using temporary local database")
		db := NewTempDB()
		return db, Fill(db)
	} else {
		logger.Printf("Using local filesystem database at %s\n", url)
		return NewLocalfs(url, logger), nil
	}
}

type SongMeta = types.SongMeta

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
					meta, err := db.NewSong(SongMeta{
						Name:     fmt.Sprintf("song%d", sng),
						Artist:   fmt.Sprintf("%cartist%d", ltr, art),
						Album:    fmt.Sprintf("album%d", alb),
						TrackNum: sng + 1,
					})
					if err != nil {
						return err
					}

					_, err = db.UpdateChords(meta.ID, []byte("sample chords go here"))
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
