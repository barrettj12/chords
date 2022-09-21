// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// backend/dblayer/tempdb.go
// A temporary (non-persistent) database stored in Go memory. Good for testing

package dblayer

import "fmt"

type song struct {
	SongMeta
	Chords
}

type TempDB struct {
	// map from id -> song
	data   map[string]*song
	nextID int
}

func (t *TempDB) GetArtists() ([]string, error) {
	artists := set[string]{}
	for _, row := range t.data {
		artists.add(row.artist)
	}
	return artists.toSlice(), nil
}

func (t *TempDB) GetSongs(artist, id string) ([]SongMeta, error) {
	songs := []SongMeta{}

	for _, s := range t.data {
		matches := true
		if artist != "" && s.artist != artist {
			matches = false
		}
		if id != "" && s.id != id {
			matches = false
		}
		if matches {
			songs = append(songs, s.SongMeta)
		}
	}

	return songs, nil
}

func (t *TempDB) NewSong(meta SongMeta) (SongMeta, error) {
	idStr := fmt.Sprint(t.nextID)
	meta.id = idStr
	t.nextID++

	t.data[idStr] = &song{meta, []byte{}}
	return meta, nil
}

func (t *TempDB) UpdateSong(id string, meta SongMeta) (SongMeta, error) {
	song, ok := t.data[id]
	if !ok {
		return SongMeta{}, songNotFound(id)
	}

	meta.id = id
	song.SongMeta = meta
	return meta, nil
}

func (t *TempDB) DeleteSong(id string) error {
	delete(t.data, id)
	return nil
}

func (t *TempDB) GetChords(id string) (Chords, error) {
	song, ok := t.data[id]
	if !ok {
		return Chords{}, songNotFound(id)
	}
	return song.Chords, nil
}

func (t *TempDB) UpdateChords(id string, chords Chords) (Chords, error) {
	song, ok := t.data[id]
	if !ok {
		return Chords{}, songNotFound(id)
	}
	song.Chords = chords
	return chords, nil
}

// Set type
type set[T comparable] map[T]struct{}

func (s *set[T]) add(t T) {
	if _, ok := (*s)[t]; !ok {
		(*s)[t] = struct{}{}
	}
}

func (s *set[T]) toSlice() []T {
	slice := []T{}
	for t := range *s {
		slice = append(slice, t)
	}
	return slice
}

// Helper functions
func songNotFound(id string) error {
	return fmt.Errorf("no song found for id %s", id)
}
