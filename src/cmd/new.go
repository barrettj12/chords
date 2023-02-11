package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"

	"github.com/barrettj12/chords/src/dblayer"
)

// Interactively add a new song to the local DB.
//
//	usage: chords new
func new(st state, args []string) {
	db := dblayer.NewLocalfs(st.dbPath, log.Default())
	s := bufio.NewScanner(os.Stdin)

	songName := promptf(s, "Song name: ")
	// Attempt to put song name in title case
	properSongName := properTitle(songName)
	if songName != properSongName {
		resp := promptf(s, "Rename to %q? [y/n]: ", properSongName)
		if resp == "y" {
			songName = properSongName
		}
	}

	id := getID(songName)
	fmt.Printf("Suggested ID: %s\n", id)
	idResp := promptf(s, "ID (enter to use default): ")
	if idResp != "" {
		id = idResp
	}

	// Check ID in DB
	for {
		_, err := db.GetChords(id)
		if err != nil {
			// This ID not already in DB - OK to continue
			break
		}

		// Song exists in DB - choose new ID
		fmt.Printf("ID %q is taken, please choose another.\n", id)
		id = promptf(s, "ID: ")
	}

	// TODO: use discogs API to get this metadata
	// https://github.com/irlndts/go-discogs
	artist := promptf(s, "Artist: ")

	// TODO: suggest existing albums for artist
	album := promptf(s, "Album: ")
	var trackNum int
	if album != "" {
		trackNum, _ = strconv.Atoi(promptf(s, "Track number: "))
	}

	_, err := db.NewSong(dblayer.SongMeta{
		ID:       id,
		Name:     songName,
		Artist:   artist,
		Album:    album,
		TrackNum: trackNum,
	})
	if err != nil {
		log.Fatalf("Error writing to DB: %s", err)
	}

	// TODO: how to do multiline prompt?
	// chords := prompt(s, "Chords: ")
	cmd := exec.Command("code", fmt.Sprintf("%s/%s/chords.txt", st.dbPath, id))
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Error creating chords: %s", err)
	}

	// TODO: wait for editor to close, then sync
}

// promptf prints the question to stdout, then reads a line from the provided
// scanner, and returns this value.
func promptf(s *bufio.Scanner, f string, v ...any) string {
	fmt.Printf(f, v...)
	s.Scan()
	return s.Text()
}

func getID(song string) string {
	id := ""
	words := strings.Fields(song)
	for _, word := range words {
		for i, c := range word {
			if !isAlphanumeric(c) {
				continue
			}
			if i == 0 {
				id += string(unicode.ToUpper(c))
			} else {
				id += string(unicode.ToLower(c))
			}
		}
	}
	return id
}

// https://www.jeremymorgan.com/tutorials/go/learn-golang-casing/
// TODO: this can't lowercase incorrectly capitalised words e.g. "The"
func properTitle(input string) string {
	words := strings.Split(input, " ")
	smallwords := " a an on the to "

	for index, word := range words {
		if strings.Contains(smallwords, " "+word+" ") && word != string(word[0]) {
			words[index] = word
		} else {
			words[index] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}

func isAlphanumeric(c rune) bool {
	return ('0' <= c && c <= '9') ||
		('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')
}
