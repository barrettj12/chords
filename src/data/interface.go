package data

import (
	"context"
	"log"

	"github.com/barrettj12/chords/src/dblayer"
)

type ChordsDBv1 interface {
	Artists(context.Context, ArtistsFilters) ([]Artist, error)
	Albums(context.Context, AlbumsFilters) ([]Album, error)
	Songs(context.Context, SongsFilters) ([]Song, error)
}

type ArtistsFilters struct {
	ID        ArtistID
	RelatedTo ArtistID
	Album     AlbumID
}

type AlbumsFilters struct {
	ID     AlbumID
	Artist ArtistID
}

type SongsFilters struct {
	ID    SongID
	Album AlbumID
}

func GetDBv1(url string, logger *log.Logger) (ChordsDBv1, error) {
	db, err := dblayer.GetDB(url, logger)
	if err != nil {
		return nil, err
	}
	return &ChordsDBv1Shim{db}, nil
}
