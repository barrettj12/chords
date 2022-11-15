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
	"log"
	"os"
	"strconv"

	"github.com/barrettj12/chords/src/dblayer"
	"github.com/barrettj12/chords/src/server"
)

func main() {
	// Try to read logging flags from LOG_FLAGS environment variable
	// Invalid/unset values will just default to 0 (no flags)
	flags, _ := strconv.Atoi(os.Getenv("LOG_FLAGS"))

	// Initialise logger
	logger := log.New(os.Stdout, "", flags)

	// Set up DB
	var db dblayer.ChordsDB
	dbURL := os.Getenv("DATABASE_URL")
	db, err := dblayer.GetDB(dbURL, logger)
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	// Read port from PORT environment variable
	port := os.Getenv("PORT")
	if _, err := strconv.Atoi(port); err != nil {
		// Set default port value
		fmt.Printf("Invalid port %q: listening on port 8080 instead\n", port)
		port = "8080"
	}

	// Read authorisation key from env
	authKey := os.Getenv("AUTH_KEY")
	if authKey == "" {
		// Try read from file
		data, err := os.ReadFile("auth_key")
		if err == nil {
			authKey = string(data)
		}
	}

	s, err := server.New(db, ":"+port, logger, authKey)
	if err != nil {
		panic(err)
	}
	err = s.Run()
	if err != nil {
		panic(err)
	}
}
