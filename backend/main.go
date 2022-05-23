// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// backend/main.go
// Entry point for the backend

package main

import (
	"github.com/barrettj12/chords/backend/dblayer"
	"github.com/barrettj12/chords/backend/server"
	"os"
	"strconv"
)

func main() {
	// Set up DB
	db := &dblayer.TempDB{}
	dblayer.Fill(db)

	// Try to read logging flags from env:CHORDS_LOG
	var flag int
	flagstr, flagset := os.LookupEnv("CHORDS_LOG")
	if flagset {
		flag, _ = strconv.Atoi(flagstr)
	}

	s := server.New(db, ":8080", flag)
	s.Start()
}
