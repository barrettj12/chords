// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// src/dblayer/localfs.go
// A database layer for a local filesystem.

package dblayer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"

	gqltypes "github.com/barrettj12/chords/gqlgen/types"
	"github.com/barrettj12/chords/src/search"
	"github.com/barrettj12/chords/src/types"
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
	// Index for text search
	index *search.Index
}

func NewLocalfs(basedir string, logger *log.Logger) *localfs {
	db := &localfs{
		basedir: basedir,
		log:     logger,
	}
	db.makeIndex()
	return db
}

// makeIndex creates a search.Index and initialises it with all the songs
// already in the database.
func (l *localfs) makeIndex() {
	index, err := search.NewIndex()
	if err != nil {
		l.log.Printf("WARNING could not create search index: %v", err)
		return
	}
	l.index = index

	var dirs []fs.DirEntry
	dirs, err = os.ReadDir(l.basedir)
	if err != nil {
		l.log.Printf("WARNING could not create search index: reading %q: %v", l.basedir, err)
		return
	}

	for _, d := range dirs {
		if !d.Type().IsDir() {
			continue
		}

		meta, err := l.getMeta(d.Name())
		if err != nil {
			l.log.Printf("WARNING creating search index: getting metadata for ID %q: %v", d.Name(), err)
			continue
		}

		l.index.Add(meta)
	}
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

		meta, err := l.getMeta(d.Name())
		if err != nil {
			l.log.Printf("WARNING getting metadata for ID %q: %v", d.Name(), err)
			continue
		}

		artists.add(meta.Artist)
	}

	return artists.toSlice(), nil
}

func (l *localfs) GetSongs(artist, id, query string) ([]SongMeta, error) {
	songs := []SongMeta{}
	dirs, err := os.ReadDir(l.basedir)
	if err != nil {
		return nil, err
	}

	var queryMatcher *regexp.Regexp
	if query != "" {
		queryMatcher, err = regexp.Compile("(?i)" + query) // case insensitive
		if err != nil {
			l.log.Printf("WARNING ignoring query %q: %v", query, err)
		}
	}

	for _, d := range dirs {
		if !d.Type().IsDir() {
			continue
		}
		if id != "" && d.Name() != id {
			continue
		}

		meta, err := l.getMeta(d.Name())
		if err != nil {
			l.log.Printf("WARNING getting metadata for ID %q: %v", d.Name(), err)
			continue
		}

		// Check if artist matches
		if artist != "" && meta.Artist != artist {
			continue
		}
		if queryMatcher != nil && !queryMatcher.MatchString(meta.Name) {
			continue
		}
		songs = append(songs, meta)
	}

	return songs, nil
}

func (l *localfs) NewSong(meta SongMeta) (SongMeta, error) {
	if !l.checkID(meta.ID) {
		return SongMeta{}, fmt.Errorf("id %q already in use", meta.ID)
	}

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

	// Update search index
	err = l.index.Add(meta)
	if err != nil {
		l.log.Printf("WARNING error updating index: %v", err)
	}

	// Create empty chords.txt
	err = os.WriteFile(filepath.Join(newDir, "chords.txt"), []byte{}, os.ModePerm)
	if err != nil {
		return SongMeta{}, err
	}

	return meta, nil
}

// checkID returns true if the given ID is available (not already in use).
func (l *localfs) checkID(id string) bool {
	_, err := os.Stat(filepath.Join(l.basedir, id))
	return errors.Is(err, os.ErrNotExist)
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

	// Update search index
	err = l.index.Remove(id)
	if err != nil {
		l.log.Printf("WARNING error updating index: %v", err)
	}
	err = l.index.Add(meta)
	if err != nil {
		l.log.Printf("WARNING error updating index: %v", err)
	}

	return meta, nil
}

func (l *localfs) DeleteSong(id string) error {
	dir := filepath.Join(l.basedir, id)
	err := os.RemoveAll(dir)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	// Update search index
	err = l.index.Remove(id)
	if err != nil {
		l.log.Printf("WARNING error updating index: %v", err)
	}

	return nil
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
		l.log.Printf("WARNING see-also.json not found")
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

func (l *localfs) Search(query string) ([]SongMeta, error) {
	ids, err := l.index.Search(query)
	if err != nil {
		return nil, err
	}

	songs := []types.SongMeta{}
	for _, id := range ids {
		meta, err := l.getMeta(id)
		if err != nil {
			l.log.Printf("WARNING getting metadata for ID %q: %v", id, err)
			continue
		}
		songs = append(songs, meta)

		// Limit to first 10 results
		if len(songs) >= 10 {
			break
		}
	}

	return songs, nil
}

func (l *localfs) getMeta(id string) (meta types.SongMeta, err error) {
	metaPath := filepath.Join(l.basedir, id, "meta.json")
	file, err := os.Open(metaPath)
	if err != nil {
		return meta, err
	}

	err = json.NewDecoder(file).Decode(&meta)
	return meta, err
}

// ArtistsV1 implements ChordsDBv1.
func (l *localfs) ArtistsV1(ctx context.Context) ([]*gqltypes.Artist, error) {
	artistNames, err := l.GetArtists()
	if err != nil {
		return nil, err
	}

	//// Generate see also map
	//seeAlsoMap := map[string][]string{}
	//
	//seeAlsoPath := filepath.Join(l.basedir, "see-also.json")
	//seeAlsoFile, err := os.Open(seeAlsoPath)
	//seeAlsos := [][]string{}
	//err = json.NewDecoder(seeAlsoFile).Decode(&seeAlsos)
	//if err != nil {
	//	l.log.Printf("WARNING couldn't unmarshal see also data: %v", err)
	//}
	//
	//for _, grp := range seeAlsos {
	//	seeAlsoMap[grp[0]] = append(seeAlsoMap[grp[0]], grp[1])
	//	seeAlsoMap[grp[1]] = append(seeAlsoMap[grp[1]], grp[0])
	//}

	// Generate returned data
	var artists []*gqltypes.Artist
	for _, name := range artistNames {
		artists = append(artists, &gqltypes.Artist{
			Name: name,
			// TODO: how to put in albums / related artists?
		})
	}
	return artists, nil
}
