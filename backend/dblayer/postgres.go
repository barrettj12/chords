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

// Postgres represents a Postgres database.
type Postgres struct {
	db *sql.DB
}

// NewPostgres creates and initialises a Postgres DB at the given URL.
func NewPostgres(url string) (*Postgres, error) {
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

	return &Postgres{db}, nil
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

func (p *Postgres) GetArtists() ([]string, error) {
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

func (p *Postgres) GetSongs(artist string) (Songs, error) {
	return Songs{}, nil
}

func (p *Postgres) GetChords(id int) (string, error) {
	return "", nil
}

func (p *Postgres) SetChords(id int, data []byte) error {
	return nil
}

func (p *Postgres) MakeChords(nc NewChords) (int, error) {
	res, err := p.db.Exec(`
INSERT INTO	chords
VALUES (DEFAULT, $1, $2, $3, $4);`,
		nc.Artist, nc.Album, nc.Song, nc.Chords)
	if err != nil {
		return 0, fmt.Errorf("Postgres.MakeChords: %w", err)
	}
	id, err := res.RowsAffected()
	return int(id), nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}
