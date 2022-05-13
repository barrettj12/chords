package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

// HELPER FUNCTIONS

// prompt prints the question to stdout, then reads a line from the provided
// scanner, and sets the provided pointer to this value.
func prompt(s *bufio.Scanner, q string, out *string) {
	fmt.Print(q)
	s.Scan()
	*out = s.Text()
}
