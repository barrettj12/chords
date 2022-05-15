// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// backend/dblayer/interface.go
// This package contains the database layer, which is responsible for
// communicating with the database, and converting between the database's
// format and Go objects.
// This file contains interfaces/code shared between different DBs

package dblayer

import (
	"fmt"
)

type ChordsDB interface {
	GetArtists() ([]string, error)
	GetSongs(string) (Songs, error)
	GetChords(int) (string, error)
	SetChords(int, []byte) error
	MakeChords(NewChords) (int, error)
}

// album maps  song -> id
type album map[string]int

// songs maps albumName -> album
type Songs map[string]album

// NewChords represents the body of a POST /chords request.
type NewChords struct {
	Artist string `json:"artist"`
	Album  string `json:"name"`
	Song   string `json:"song"`
	Chords string `json:"chords"`
}

// DB ERRORS
type fmtErr struct {
	format string
	args   []interface{}
}

func (e fmtErr) Error() string {
	return fmt.Sprintf(e.format, e.args...)
}

func Errorf(format string, args ...interface{}) error {
	return fmtErr{format, args}
}
