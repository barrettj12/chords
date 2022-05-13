// a temporary (non-persistent) database stored in Go memory
// good for testing

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
		return "", Errorf("negative id %q invalid", id)
	} else if id >= len(*t) {
		return "", Errorf("no chords found for id %q", id)
	}

	return (*t)[id].Chords, nil
}

func (t *TempDB) SetChords(id int, data []byte) error {
	if id < 0 {
		return Errorf("negative id %q invalid", id)
	} else if id >= len(*t) {
		return Errorf("no chords found for id %q", id)
	}

	(*t)[id].Chords = string(data)
	return nil
}

func (t *TempDB) MakeChords(nc NewChords) (int, error) {
	id := len(*t)
	*t = append(*t, nc)
	return id, nil
}
