// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// src/cmd/chords.go
// The command-line interface, used to update the chords database.

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"

	"github.com/barrettj12/chords/src/client"
	"github.com/barrettj12/chords/src/dblayer"
	"github.com/barrettj12/chords/src/util"
)

func main() {
	st := initState()
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "albums":
		albums(st, args)
	case "backup":
		backup(args)
	case "count":
		count(st, args)
	case "delete", "rm", "remove":
		delete(st, args)
	case "diff":
		diff(st, args)
	case "edit":
		edit(st, args)
	case "new":
		new(st, args)
	case "pull":
		pull(args)
	case "sync":
		sync(st, args)
	case "update-chords":
		updateChords(st, args)
	case "validate":
		validate(st, args)
	default:
		fmt.Printf("unknown command %q\n", cmd)
		os.Exit(1)
	}
}

// Global state, passed to subcommands
type state struct {
	dbPath    string
	serverURL string
	authKey   string
}

func initState() state {
	dbPath, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		dbPath = "./data"
	}

	serverURL, ok := os.LookupEnv("SERVER_URL")
	if !ok {
		serverURL = "https://chords.fly.dev"
	}

	authKey, err := os.ReadFile("auth_key")
	if err != nil {
		fmt.Printf("WARNING: couldn't read auth_key: %v\n", err)
	} else {
		fmt.Println("INFO: using auth key from file")
	}

	return state{
		dbPath,
		serverURL,
		string(authKey),
	}
}

// pull gets chords from server.
func pull(args []string) {
	fmt.Println("pulling", args)
}

// backup makes a full backup of the database.
func backup(args []string) {
	fmt.Println("backing up", args)
}

// sync copies chords from local db to remote
//
//	sync [song-ids...]
func sync(st state, args []string) {
	db := dblayer.NewLocalfs(st.dbPath, log.Default())
	c, err := client.NewClient(st.serverURL, st.authKey)
	check(err)

	ids := args
	var songs []dblayer.SongMeta
	if len(ids) == 0 {
		// Sync all songs
		var err error
		songs, err = db.GetSongs("", "", "")
		check(err)
	} else {
		songs = make([]dblayer.SongMeta, 0, len(ids))
		for _, id := range ids {
			dbSongs, err := db.GetSongs("", id, "")
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
		songs, err := c.GetSongs(nil, &localSong.ID, nil)
		check(err)
		if len(songs) == 0 {
			// Song doesn't exist in remote DB
			_, err := c.NewSong(localSong)
			check(err)
		} else if songs[0] != localSong {
			// Update song in remote DB
			_, err := c.UpdateSong(localSong.ID, localSong)
			check(err)
		}

		// Sync chords
		chords, err := db.GetChords(localSong.ID)
		check(err)
		c.UpdateChords(localSong.ID, chords)
		check(err)
	}
}

// Update chords via the POST /api/v0/chords endpoint
//
//	chords update-chords <song-id> [path-to-chords-file]
func updateChords(st state, args []string) {
	c, err := client.NewClient(st.serverURL, st.authKey)
	check(err)
	songID := args[0]
	var path string
	if len(args) > 1 {
		path = args[1]
	} else {
		path = fmt.Sprintf("./data/%s/chords.txt", songID)
	}

	chords, err := os.ReadFile(path)
	check(err)

	_, err = c.UpdateChords(songID, chords)
	check(err)
}

// Count number of songs for each artist
//
//	count [artists...]
func count(st state, args []string) {
	artists := args
	c, err := client.NewClient(st.serverURL, st.authKey)
	check(err)

	counts := make(map[string]int, len(artists))
	songs, err := c.GetSongs(nil, nil, nil)
	check(err)

	for _, song := range songs {
		counts[song.Artist]++
	}

	if len(artists) == 0 {
		for artist := range counts {
			artists = append(artists, artist)
		}
	}

	// Sort artists in descending order
	sort.Slice(artists, func(i, j int) bool {
		ci := counts[artists[i]]
		cj := counts[artists[j]]

		if ci == cj {
			return util.LessTitle(artists[i], artists[j])
		}
		return ci > cj
	})

	for _, artist := range artists {
		numSongs := counts[artist]
		if numSongs > 0 {
			fmt.Printf("%d\t%s\n", numSongs, artist)
		}
	}
}

// List albums and their tracks
//
//	usage: chords albums [artists...]
func albums(st state, args []string) {
	c, err := client.NewClient(st.serverURL, st.authKey)
	check(err)

	artists := args

	songs := []dblayer.SongMeta{}
	if len(artists) == 0 {
		// Get songs for all artists
		songs, err = c.GetSongs(nil, nil, nil)
		check(err)
	} else {
		for _, artist := range artists {
			newSongs, err := c.GetSongs(&artist, nil, nil)
			check(err)
			songs = append(songs, newSongs...)
		}
	}

	//            artist -> album -> (trackNum, song)
	type songWithNum struct {
		name string
		num  int
	}
	albums := map[string]map[string][]songWithNum{}

	for _, song := range songs {
		if albums[song.Artist] == nil {
			albums[song.Artist] = map[string][]songWithNum{}
		}
		albums[song.Artist][song.Album] = append(
			albums[song.Artist][song.Album], songWithNum{song.Name, song.TrackNum})
	}

	// Print albums
	for artist, albumMap := range albums {
		fmt.Println(artist)
		for album, tracks := range albumMap {
			if album == "" {
				fmt.Println("  (no album)")
				for _, song := range tracks {
					fmt.Printf("    %s\n", song.name)
				}

			} else {
				fmt.Printf("  %s\n", album)
				sort.Slice(tracks, func(i, j int) bool {
					return tracks[i].num < tracks[j].num
				})

				for _, song := range tracks {
					fmt.Printf("    %d. %s\n", song.num, song.name)
				}
			}
		}
	}
}

// Open chords for editing
//
//	chords edit <id>
func edit(st state, args []string) {
	id := args[0]
	// Check if ID exists in DB
	db := dblayer.NewLocalfs(st.dbPath, log.Default())
	_, err := db.GetChords(id)
	if err != nil {
		// This ID not already in DB - OK to continue
		log.Fatalf("no chords found with ID %q", id)
	}

	editorCmd := exec.Command("code", "--wait", fmt.Sprintf("%s/%s/chords.txt", st.dbPath, id))
	err = editorCmd.Start()
	if err != nil {
		log.Fatalf("Error opening chords: %s", err)
	}

	// Wait for editor to close, then sync
	fmt.Println("Waiting for editor to close")
	err = editorCmd.Wait()
	if err != nil {
		log.Fatalf("Error editing chords: %s", err)
	}
	sync(st, []string{id})
}

// Delete a set of chords locally and remotely
//
//	chords delete <id>
func delete(st state, args []string) {
	id := args[0]

	c, err := client.NewClient(st.serverURL, st.authKey)
	check(err)

	c.DeleteSong(id)
	check(err)

	// TODO: delete locally
}

// HELPER FUNCTIONS

func check(err error) {
	if err != nil {
		panic(err)
	}
}
