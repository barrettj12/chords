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
