package dblayer

import "io"

type ChordsDB interface {
	GetArtists() ([]string, error)
	GetSongs(string) (Songs, error)
	GetChords(int) (string, error)
	SetChords(int, io.ReadCloser) error
	MakeChords(io.ReadCloser) (int, error)
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
