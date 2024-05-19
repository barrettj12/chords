// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// backend/dblayer/tempdb.go
// A temporary (non-persistent) database stored in Go memory. Good for testing

package dblayer

import (
	"fmt"
	"github.com/barrettj12/chords/src/types"
)

type song struct {
	SongMeta
	Chords
}

type tempDB struct {
	// map from id -> song
	data   map[string]*song
	nextID int
}

// Return a correctly initialised tempDB.
func NewTempDB() *tempDB {
	return &tempDB{
		make(map[string]*song),
		1,
	}
}

func (t *tempDB) GetArtists() ([]string, error) {
	artists := set[string]{}
	for _, row := range t.data {
		artists.add(row.Artist)
	}
	return artists.toSlice(), nil
}

func (t *tempDB) GetSongs(artist, id, query string) ([]SongMeta, error) {
	songs := []SongMeta{}

	for _, s := range t.data {
		if artist != "" && s.Artist != artist {
			continue
		}
		if id != "" && s.ID != id {
			continue
		}
		// TODO: handle `query` param
		songs = append(songs, s.SongMeta)
	}

	return songs, nil
}

func (t *tempDB) NewSong(meta SongMeta) (SongMeta, error) {
	idStr := fmt.Sprint(t.nextID)
	meta.ID = idStr
	t.nextID++

	t.data[idStr] = &song{meta, []byte{}}
	return meta, nil
}

func (t *tempDB) UpdateSong(id string, meta SongMeta) (SongMeta, error) {
	song, ok := t.data[id]
	if !ok {
		return SongMeta{}, songNotFound(id)
	}

	meta.ID = id
	song.SongMeta = meta
	return meta, nil
}

func (t *tempDB) DeleteSong(id string) error {
	delete(t.data, id)
	return nil
}

func (t *tempDB) GetChords(id string) (Chords, error) {
	song, ok := t.data[id]
	if !ok {
		return Chords{}, songNotFound(id)
	}
	return song.Chords, nil
}

func (t *tempDB) UpdateChords(id string, chords Chords) (Chords, error) {
	song, ok := t.data[id]
	if !ok {
		return Chords{}, songNotFound(id)
	}
	song.Chords = chords
	return chords, nil
}

func (t *tempDB) SeeAlso(artist string) ([]string, error) {
	// TODO: fill this in
	return nil, nil
}

func (t *tempDB) Search(query string) ([]types.SearchResult, error) {
	// TODO: fill this in
	return nil, nil
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
