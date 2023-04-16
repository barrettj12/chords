// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// backend/dblayer/postgres.go
// A wrapper for a Postgres database.

package dblayer

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// postgres represents a postgres database.
type postgres struct {
	db *sql.DB
}

// NewPostgres creates and initialises a Postgres DB at the given URL.
func NewPostgres(url string) (*postgres, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Ensure table `chords` exists
	// TODO
	err = initDB(db)
	if err != nil {
		return nil, err
	}

	return &postgres{db}, nil
}

// initDB creates the required tables in the given database.
func initDB(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE chords (
	id       SERIAL PRIMARY KEY,
	artist   TEXT,
	album    TEXT,
	song     TEXT,
	data     TEXT
);`)
	return err
}

func (p *postgres) GetArtists() ([]string, error) {
	rows, err := p.db.Query(`
SELECT DISTINCT artist
FROM chords
`)
	if err != nil {
		return nil, fmt.Errorf("Postgres.GetArtists: %w", err)
	}
	defer rows.Close()

	// Turn `rows` into a []string
	artists := []string{}
	for rows.Next() {
		var artist string
		rows.Scan(&artist)
		if err != nil {
			return nil, fmt.Errorf("Postgres.GetArtists: %w", err)
		}
		artists = append(artists, artist)
	}

	return artists, nil
}

func (p *postgres) GetSongs(artist, id string) ([]SongMeta, error) {
	// TODO: fill this in
	return nil, nil
}

func (p *postgres) NewSong(SongMeta) (SongMeta, error) {
	// TODO: fill this in
	return SongMeta{}, nil
	// 	res, err := p.db.Exec(`
	// INSERT INTO	chords
	// VALUES (DEFAULT, $1, $2, $3, $4);`,
	// 		nc.Artist, nc.Album, nc.Song, nc.Chords)
	// 	if err != nil {
	// 		return 0, fmt.Errorf("Postgres.MakeChords: %w", err)
	// 	}
	// 	id, err := res.RowsAffected()
	// 	return int(id), nil
}

func (p *postgres) UpdateSong(id string, meta SongMeta) (SongMeta, error) {
	// TODO: fill this in
	return SongMeta{}, nil
}

func (p *postgres) DeleteSong(id string) error {
	// TODO: fill this in
	return nil
}

func (p *postgres) GetChords(id string) (Chords, error) {
	// TODO: fill this in
	return Chords{}, nil
}

func (p *postgres) UpdateChords(id string, chords Chords) (Chords, error) {
	// TODO: fill this in
	return Chords{}, nil
	// 	_, err := p.db.Exec(`
	// UPDATE chords
	// SET data = $1
	// WHERE id = $2;`,
	// 		string(data), id)
	// 	return err
}

func (p *postgres) SeeAlso(artist string) ([]string, error) {
	// TODO: fill this in
	return nil, nil
}
