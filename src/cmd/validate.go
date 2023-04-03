package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/barrettj12/chords/src/dblayer"
)

// Validate local database
//
//	usage: chords validate
//
// TODO:
// - extra checks
//   - no album has same track num twice
//   - similarity of artist/album names
//
// - option to set the "max track num" above which we will warn
// - option to silence some checks by name
func validate(st state, args []string) {
	datadir := st.dbPath
	entries, err := os.ReadDir(datadir)
	check(err)

	for _, entry := range entries {
		path := filepath.Join(st.dbPath, entry.Name())

		// Check it's a directory
		if !entry.IsDir() {
			log.Printf("WARNING: %q is not a directory\n", path)
			continue
		}

		// Read the dir - should contain exactly two files
		//   meta.json
		//   chords.txt
		files, err := os.ReadDir(path)
		if err != nil {
			log.Printf("WARNING: couldn't read dir %q: %v\n", path, err)
			continue
		}

		metaFound := false
		chordsFound := false
		for _, file := range files {
			fpath := filepath.Join(path, file.Name())

			// Check it's a plain file
			if file.IsDir() {
				log.Printf("WARNING: %q is a directory\n", fpath)
				continue
			}

			// Check name
			switch file.Name() {
			case "meta.json":
				metaFound = true
				validateMeta(fpath, entry.Name())
			case "chords.txt":
				chordsFound = true
				validateChords(fpath)
			default:
				log.Printf("WARNING: unexpected file found: %q\n", fpath)
			}
		}

		if !metaFound {
			log.Printf("WARNING: %q: no meta.json found\n", path)
		}
		if !chordsFound {
			log.Printf("WARNING: %q: no chords.txt found\n", path)
		}
	}
}

func validateMeta(fpath, dirName string) {
	data, err := os.ReadFile(fpath)
	if err != nil {
		log.Printf("WARNING: couldn't read %q: %v\n", fpath, err)
		return
	}

	song := dblayer.SongMeta{}
	err = json.Unmarshal(data, &song)
	if err != nil {
		log.Printf("WARNING: couldn't parse %q: %v\n", fpath, err)
		return
	}

	// ID should match dir
	if song.ID != dirName {
		log.Printf("WARNING: %q: id %q doesn't match dir name %q",
			fpath, song.ID, dirName)
	}

	// Check all fields are defined
	if song.Name == "" {
		log.Printf("WARNING: %q: song name is empty\n", fpath)
	}
	if song.Artist == "" {
		log.Printf("WARNING: %q: artist is empty\n", fpath)
	}
	// TODO: we allow album == "", but only if album/trackNum are not present
	// in meta.json
	if song.Album == "" {
		log.Printf("WARNING: %q: album is empty\n", fpath)
	}
	if song.TrackNum <= 0 {
		log.Printf("WARNING: %q: trackNum should be at least 1\n", fpath)
	}
	if song.TrackNum > 20 {
		log.Printf("WARNING: %q: trackNum might be too large?\n", fpath)
	}

	// TODO: check json fmt with jq
}

func validateChords(fpath string) {
	data, err := os.ReadFile(fpath)
	if err != nil {
		log.Printf("WARNING: couldn't read %q: %v\n", fpath, err)
		return
	}

	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		log.Printf("WARNING: %q: chords are empty\n", fpath)
	}
}
