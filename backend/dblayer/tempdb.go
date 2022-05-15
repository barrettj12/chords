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
	artists := []string{}
	for _, row := range *t {
		artists = append(artists, row.Artist)
	}
	return artists, nil
}

func (t *TempDB) GetSongs(artist string) (Songs, error) {
	songs := Songs{}
	for i, row := range *t {
		if artist == row.Artist {
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
