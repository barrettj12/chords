package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/barrettj12/chords/src/dblayer"
)

// Interactively add a new song to the local DB.
//
//	usage: chords new
func new(st state, args []string) {
	db := dblayer.NewLocalfs(st.dbPath, log.Default())
	s := bufio.NewScanner(os.Stdin)

	songName := prompt(s, "Song name: ")
	id := getID(songName)
	fmt.Printf("Suggested ID: %s\n", id)

	idResp := prompt(s, "ID (enter to use default): ")
	if idResp != "" {
		id = idResp
	}

	// TODO: check ID in DB

	// TODO: use discogs API to get this metadata
	// https://github.com/irlndts/go-discogs
	artist := prompt(s, "Artist: ")
	album := prompt(s, "Album: ")
	trackNum, _ := strconv.Atoi(prompt(s, "Track number: "))

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
	cmd := exec.Command("gedit", fmt.Sprintf("%s/%s/chords.txt", st.dbPath, id))
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error creating chords: %s", err)
	}
}

// prompt prints the question to stdout, then reads a line from the provided
// scanner, and returns this value.
func prompt(s *bufio.Scanner, q string) string {
	fmt.Print(q)
	s.Scan()
	return s.Text()
}

func getID(song string) string {
	id := ""
	for _, c := range properTitle(song) {
		if isAlphanumeric(c) {
			id += string(c)
		}
	}
	return id
}

// https://www.jeremymorgan.com/tutorials/go/learn-golang-casing/
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
