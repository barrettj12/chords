package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/barrettj12/chords/src/dblayer"
	c "github.com/barrettj12/collections"
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

	err := os.MkdirAll(filepath.Join(st.dbPath, id), 0777)
	if err != nil {
		log.Fatalf("Error creating chords directory: %s", err)
	}

	editor := chooseEditor()
	editorCmd := exec.Command(editor, "--wait", filepath.Join(st.dbPath, id, "chords.txt"))
	err = editorCmd.Start()
	if err != nil {
		log.Fatalf("Error creating chords: %s", err)
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

	// Create meta.json
	data, err := json.Marshal(dblayer.SongMeta{
		ID:       id,
		Name:     songName,
		Artist:   artist,
		Album:    album,
		TrackNum: trackNum,
	})
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}
	err = os.MkdirAll(filepath.Join(st.dbPath, id), 0755)
	if err != nil {
		log.Fatalf("Error creating folder: %v", err)
	}
	err = os.WriteFile(filepath.Join(st.dbPath, id, "meta.json"), data, os.ModePerm)
	if err != nil {
		log.Fatalf("Error writing metadata: %v", err)
	}

	// Wait for editor to close, then sync
	fmt.Println("Waiting for editor to close")
	err = editorCmd.Wait()
	if err != nil {
		log.Fatalf("Error editing chords: %s", err)
	}
	sync(st, []string{id})
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

func isAlphanumeric(c rune) bool {
	return ('0' <= c && c <= '9') ||
		('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')
}

func properTitle(input string) string {
	words := strings.Split(input, " ")

	// These words shouldn't be capitalised in titles
	smallWords := c.AsSet([]string{"a", "an", "and", "for", "in", "of", "on", "or", "the", "to"})

	for i, word := range words {
		// First word should always be capitalised
		if i == 0 || !smallWords.Contains(strings.TrimRight(word, ",")) {
			words[i] = capitalise(word)
		}
	}

	return strings.Join(words, " ")
}

func capitalise(word string) string {
	var capWord string
	for i, c := range word {
		if i == 0 {
			capWord += string(unicode.ToUpper(c))
		} else {
			capWord += string(c)
		}
	}
	return capWord
}

// Choose an editor for the chords
func chooseEditor() string {
	editor := os.Getenv("CHORDS_EDITOR")
	if editor != "" {
		// Check the specified editor exists on path
		path, err := exec.LookPath(editor)
		if err != nil {
			return path
		}
		fmt.Printf("WARNING: %q not found on path\n", editor)
	}

	// Default to gedit
	return "gedit"
}
