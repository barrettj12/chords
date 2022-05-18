// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// backend/dblayer/tempdb.go
// A temporary (non-persistent) database stored in Go memory. Good for testing

package dblayer

type TempDB []NewChords

func (t *TempDB) GetArtists() ([]string, error) {
	artists := set[string]{}
	for _, row := range *t {
		artists.add(row.Artist)
	}
	return artists.toSlice(), nil
}

func (t *TempDB) GetSongs(artist string) (Songs, error) {
	songs := make(Songs)
	for i, row := range *t {
		if artist == row.Artist {
			if songs[row.Album] == nil {
				songs[row.Album] = make(album)
			}
			songs[row.Album][row.Song] = i
		}
	}
	return songs, nil
}

func (t *TempDB) GetChords(id int) (string, error) {
	if id < 0 {
		return "", Errorf("negative id %d invalid", id)
	} else if id >= len(*t) {
		return "", Errorf("no chords found for id %d", id)
	}

	return (*t)[id].Chords, nil
}

func (t *TempDB) SetChords(id int, data []byte) error {
	if id < 0 {
		return Errorf("negative id %d invalid", id)
	} else if id >= len(*t) {
		return Errorf("no chords found for id %d", id)
	}

	(*t)[id].Chords = string(data)
	return nil
}

func (t *TempDB) MakeChords(nc NewChords) (int, error) {
	id := len(*t)
	*t = append(*t, nc)
	return id, nil
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
