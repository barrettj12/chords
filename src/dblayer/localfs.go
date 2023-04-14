// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// src/dblayer/localfs.go
// A database layer for a local filesystem.

package dblayer

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

// Store data in an attached filesystem
// The file structure is as follows:
//   basedir
//   ├─ [id1]
//   │  ├─ meta.json
//   │  └─ chords.txt
//   ├─ [id2]
//   │  ├─ meta.json
//   │  └─ chords.txt
//   ...

type localfs struct {
	basedir string
	log     *log.Logger
}

func NewLocalfs(basedir string, logger *log.Logger) ChordsDB {
	return &localfs{basedir, logger}
}

func (l *localfs) GetArtists() ([]string, error) {
	artists := set[string]{}
	dirs, err := os.ReadDir(l.basedir)
	if err != nil {
		return nil, err
	}

	for _, d := range dirs {
		if !d.Type().IsDir() {
			continue
		}

		metaPath := filepath.Join(l.basedir, d.Name(), "meta.json")
		data, err := os.ReadFile(metaPath)
		if err != nil {
			log.Printf("WARNING reading file %s: %v", metaPath, err)
			continue
		}

		meta := &SongMeta{}
		err = json.Unmarshal(data, meta)
		if err != nil {
			log.Printf("WARNING unmarshaling json: %v", err)
			continue
		}
		artists.add(meta.Artist)
	}

	return artists.toSlice(), nil
}

func (l *localfs) GetSongs(artist, id string) ([]SongMeta, error) {
	songs := []SongMeta{}
	dirs, err := os.ReadDir(l.basedir)
	if err != nil {
		return nil, err
	}

	for _, d := range dirs {
		if !d.Type().IsDir() {
			continue
		}
		if id != "" && d.Name() != id {
			continue
		}

		metaPath := filepath.Join(l.basedir, d.Name(), "meta.json")
		data, err := os.ReadFile(metaPath)
		if err != nil {
			log.Printf("WARNING reading file %s: %v", metaPath, err)
			continue
		}

		meta := &SongMeta{}
		err = json.Unmarshal(data, meta)
		if err != nil {
			log.Printf("WARNING unmarshaling json: %v", err)
			continue
		}

		// Check if artist matches
		if artist != "" && meta.Artist != artist {
			continue
		}
		songs = append(songs, *meta)
	}

	return songs, nil
}

func (l *localfs) NewSong(meta SongMeta) (SongMeta, error) {
	meta.ID = l.checkID(meta.ID)

	// Create new dir
	newDir := filepath.Join(l.basedir, meta.ID)
	err := os.Mkdir(newDir, os.ModePerm)
	if err != nil {
		return SongMeta{}, err
	}

	// Create meta.json
	data, err := json.Marshal(meta)
	if err != nil {
		return SongMeta{}, err
	}
	err = os.WriteFile(filepath.Join(newDir, "meta.json"), data, os.ModePerm)
	if err != nil {
		return SongMeta{}, err
	}

	// Create empty chords.txt
	err = os.WriteFile(filepath.Join(newDir, "chords.txt"), []byte{}, os.ModePerm)
	if err != nil {
		return SongMeta{}, err
	}

	return meta, nil
}

// Checks if the give ID is in use already, and returns it if not.
// If the ID is use, return a fresh ID.
func (l *localfs) checkID(id string) string {
	if _, err := os.Stat(filepath.Join(l.basedir, id)); errors.Is(err, os.ErrNotExist) {
		return id
	}
	return l.getNewID()
}

func (l *localfs) getNewID() string {
	for {
		b := make([]byte, 4)
		rand.Read(b)
		id := hex.EncodeToString(b)

		// Check this ID is not in use
		if _, err := os.Stat(filepath.Join(l.basedir, id)); errors.Is(err, os.ErrNotExist) {
			return id
		}
	}
}

func (l *localfs) UpdateSong(id string, meta SongMeta) (SongMeta, error) {
	dir := filepath.Join(l.basedir, id)
	// Check id exists in database
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		return SongMeta{}, err
	}

	meta.ID = id
	data, err := json.Marshal(meta)
	if err != nil {
		return SongMeta{}, err
	}
	err = os.WriteFile(filepath.Join(dir, "meta.json"), data, os.ModePerm)
	if err != nil {
		return SongMeta{}, err
	}

	return meta, nil
}

func (l *localfs) DeleteSong(id string) error {
	dir := filepath.Join(l.basedir, id)
	err := os.RemoveAll(dir)

	if errors.Is(err, os.ErrNotExist) {
		return nil // no-op
	} else {
		return err
	}
}

func (l *localfs) GetChords(id string) (Chords, error) {
	path := filepath.Join(l.basedir, id, "chords.txt")
	return os.ReadFile(path)
}

func (l *localfs) UpdateChords(id string, chords Chords) (Chords, error) {
	path := filepath.Join(l.basedir, id, "chords.txt")
	os.WriteFile(path, chords, os.ModePerm)
	return os.ReadFile(path)
}

func (l *localfs) SeeAlso(artist string) ([]string, error) {
	path := filepath.Join(l.basedir, "see-also.json")
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		// File doesn't exist - no see also data to report
		log.Printf("WARNING see-also.json not found")
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	seeAlsos := &[][]string{}
	err = json.Unmarshal(data, seeAlsos)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal see also data: %w", err)
	}

	var artists []string
	for _, grp := range *seeAlsos {
		if grp[0] == artist {
			artists = append(artists, grp[1])
		}
		if grp[1] == artist {
			artists = append(artists, grp[0])
		}
	}

	return artists, nil
}
