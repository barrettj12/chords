package dblayer

import "io"

type ChordsDB interface {
	GetArtists() ([]string, error)
	GetSongs(string) (Songs, error)
	GetChords(int) (string, error)
	SetChords(int, io.ReadCloser) error
}

// album maps  song -> id
type album map[string]int

// songs maps albumname -> album
type Songs map[string]album
