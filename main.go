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
	"strings"

	"github.com/barrettj12/chords/backend/dblayer"
	"github.com/barrettj12/chords/backend/server"
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
	if strings.HasPrefix(dbURL, "postgres") {
		fmt.Printf("Using Postgres database at %s\n", dbURL)
		var err error
		db, err = dblayer.NewPostgres(dbURL)
		if err != nil {
			panic(err)
		}
	} else if dbURL == "" {
		fmt.Println("Using temporary local database")
		db = dblayer.NewTempDB()
		err := dblayer.Fill(db)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("Using local filesystem database at %s\n", dbURL)
		db = dblayer.NewLocalfs(dbURL, logger)
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

	s := server.New(db, ":"+port, logger, authKey)
	s.Start()
}
