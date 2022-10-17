// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// cmd/chords.go
// The command-line interface used to update the chords database.

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/barrettj12/chords/backend/dblayer"
	"github.com/barrettj12/chords/client"
)

var SERVER_URL = "http://localhost:8080"

func main() {
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
//	sync <path/to/local/db> <server/url>
func sync(args []string) {
	db := dblayer.NewLocalfs(args[0], log.Default())
	serverURL := args[1]
	authKey, err := os.ReadFile("auth_key")
	if err != nil {
		fmt.Printf("WARNING: couldn't read auth_key: %v\n", err)
	}
	client, err := client.NewClient(serverURL, string(authKey))
	if err != nil {
		panic(err)
	}

	songs, err := db.GetSongs("", "")
	if err != nil {
		panic(err)
	}

	for _, localSong := range songs {
		songs, err := client.GetSongs(nil, &localSong.ID, nil)
		if err != nil {
			panic(err)
		}
		if len(songs) == 0 {
			// Song doesn't exist in remote DB
			_, err := client.NewSong(localSong)
			if err != nil {
				panic(err)
			}
		} else {
			// Update song in remote DB
			_, err := client.UpdateSong(localSong.ID, localSong)
			if err != nil {
				panic(err)
			}
		}

		// Sync chords
		chords, err := db.GetChords(localSong.ID)
		if err != nil {
			panic(err)
		}
		client.UpdateChords(localSong.ID, chords)
		if err != nil {
			panic(err)
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
