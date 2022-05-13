package dblayer

type ChordsDB interface {
	GetArtists() []string
	GetSongs(string) Songs
	GetChords(int) string
}

// album maps  song -> id
type album map[string]int

// songs maps albumname -> album
type Songs map[string]album
