// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// src/cmd/chords.go
// The command-line interface, used to update the chords database.

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/barrettj12/chords/src/client"
	"github.com/barrettj12/chords/src/dblayer"
)

var SERVER_URL = "http://localhost:8080"

func main() {
	// TODO: set server-url, local-db, auth-key via env variables
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "pull":
		pull(args)
	case "push":
		push(args)
	case "backup":
		backup(args)
	case "sync":
		sync(args)
	case "update-chords":
		updateChords(args)
	case "count":
		count(args)
	case "validate":
		validate(args)
	// TODO: add a command to interactively add chords to DB
	// using discogs API to get metadata
	// https://github.com/irlndts/go-discogs
	default:
		fmt.Printf("unknown command %q\n", cmd)
		os.Exit(1)
	}
}

// pull gets chords from server.
func pull(args []string) {
	fmt.Println("pulling", args)
}

// push updates chords on the server.
func push(args []string) {
	s := bufio.NewScanner(os.Stdin)
	var artist, album, song, chords string

	prompt(s, "Artist: ", &artist)
	prompt(s, "Album: ", &album)
	prompt(s, "Song: ", &song)
	prompt(s, "Chords: ", &chords)

	buf, err := json.Marshal(map[string]string{
		"artist": artist,
		"album":  album,
		"song":   song,
		"chords": chords,
	})
	if err != nil {
		log.Fatalf("Error encoding to JSON: %s", err)
	}

	http.Post(SERVER_URL+"/chords", "application/json", bytes.NewReader(buf))
}

// backup makes a full backup of the database.
func backup(args []string) {
	fmt.Println("backing up", args)
}

// sync copies chords from local db to remote
//
//	sync <path/to/local/db> <server/url> [song-ids...]
func sync(args []string) {
	db := dblayer.NewLocalfs(args[0], log.Default())
	serverURL := args[1]
	client := getClient(serverURL)

	ids := args[2:]
	var songs []dblayer.SongMeta
	if len(ids) == 0 {
		// Sync all songs
		var err error
		songs, err = db.GetSongs("", "")
		check(err)
	} else {
		songs = make([]dblayer.SongMeta, 0, len(ids))
		for _, id := range ids {
			dbSongs, err := db.GetSongs("", id)
			check(err)
			if len(dbSongs) == 0 {
				fmt.Printf("song %q not found\n", id)
			} else {
				songs = append(songs, dbSongs[0])
			}
		}
	}

	if len(songs) == 0 {
		fmt.Println("no songs to sync")
	}

	for _, localSong := range songs {
		// TODO: parallelise here using goroutines
		fmt.Printf("syncing %q\n", localSong.Name)
		songs, err := client.GetSongs(nil, &localSong.ID, nil)
		check(err)
		if len(songs) == 0 {
			// Song doesn't exist in remote DB
			_, err := client.NewSong(localSong)
			check(err)
		} else if songs[0] != localSong {
			// Update song in remote DB
			_, err := client.UpdateSong(localSong.ID, localSong)
			check(err)
		}

		// Sync chords
		chords, err := db.GetChords(localSong.ID)
		check(err)
		client.UpdateChords(localSong.ID, chords)
		check(err)
	}
}

// Update chords via the POST /api/v0/chords endpoint
//
//	chords update-chords <server-url> <song-id> [path-to-chords-file]
func updateChords(args []string) {
	serverURL := args[0]
	client := getClient(serverURL)
	songID := args[1]
	var path string
	if len(args) > 2 {
		path = args[2]
	} else {
		path = fmt.Sprintf("./data/%s/chords.txt", songID)
	}

	chords, err := os.ReadFile(path)
	check(err)

	_, err = client.UpdateChords(songID, chords)
	check(err)
}

// Count number of songs for each artist
//
//	count <server-url> [artists...]
func count(args []string) {
	serverURL := args[0]
	artists := args[1:]
	c, err := client.NewClient(serverURL, "")
	check(err)

	counts := make(map[string]int, len(artists))
	songs, err := c.GetSongs(nil, nil, nil)
	check(err)

	for _, song := range songs {
		counts[song.Artist]++
	}

	for artist, numSongs := range counts {
		if len(artists) == 0 || sliceContains(artists, artist) {
			fmt.Printf("%d\t%s\n", numSongs, artist)
		}
	}
}

// HELPER FUNCTIONS

// prompt prints the question to stdout, then reads a line from the provided
// scanner, and sets the provided pointer to this value.
func prompt(s *bufio.Scanner, q string, out *string) {
	fmt.Print(q)
	s.Scan()
	*out = s.Text()
}

func getClient(serverURL string) *client.Client {
	authKey, err := os.ReadFile("auth_key")
	if err != nil {
		fmt.Printf("WARNING: couldn't read auth_key: %v\n", err)
	} else {
		fmt.Println("using auth key from file")
	}
	client, err := client.NewClient(serverURL, string(authKey))
	check(err)
	return client
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func sliceContains[T comparable](slice []T, t T) bool {
	for _, u := range slice {
		if t == u {
			return true
		}
	}
	return false
}
