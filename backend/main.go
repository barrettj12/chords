// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// backend/main.go
// Entry point for the backend

package main

import (
	"fmt"
	"github.com/barrettj12/chords/backend/dblayer"
	"github.com/barrettj12/chords/backend/server"
	"os"
	"strconv"
)

func main() {
	// Set up DB
	db := &dblayer.TempDB{}
	_ = dblayer.Fill(db)

	// Read port from PORT environment variable
	port := os.Getenv("PORT")
	if _, err := strconv.Atoi(port); err != nil {
		// Set default port value
		fmt.Printf("Invalid port %q: listening on port 8080 instead\n", port)
		port = "8080"
	}

	// Try to read logging flags from LOG_FLAGS environment variable
	// Invalid/unset values will just default to 0 (no flags)
	flags, _ := strconv.Atoi(os.Getenv("LOG_FLAGS"))

	s := server.New(db, ":"+port, flags)
	s.Start()
}
