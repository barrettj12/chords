package data

import (
	"context"

	"github.com/barrettj12/chords/src/dblayer"
	"github.com/barrettj12/chords/src/util"
	"github.com/barrettj12/collections"
)

// ChordsDBv1Shim is an adapter from a ChordsDB to a ChordsDBv1.
type ChordsDBv1Shim struct {
	db dblayer.ChordsDB
}

func (db *ChordsDBv1Shim) Artists(_ context.Context, filters ArtistsFilters) ([]Artist, error) {
	artistNames, err := db.db.GetArtists()
	if err != nil {
		return nil, err
	}

	// Get album info from song metadata
	songs, err := db.db.GetSongs("", "", "")
	if err != nil {
		return nil, err
	}

	albums := map[string]*collections.Set[AlbumID]{}
	for _, song := range songs {
		if song.Artist == "" {
			continue
		}
		if albums[song.Artist] == nil {
			albums[song.Artist] = collections.NewSet[AlbumID](0)
		}
		albums[song.Artist].Add(MakeAlbumID(song.Album))
	}

	artists := []Artist{}
	for _, artistName := range artistNames {
		// Get "see also" data to fill related artists
		seeAlso, _ := db.db.SeeAlso(artistName)
		relatedArtists := []ArtistID{}
		for _, relatedArtist := range seeAlso {
			relatedArtists = append(relatedArtists, MakeArtistID(relatedArtist))
		}

		artists = append(artists, Artist{
			ID:             MakeArtistID(artistName),
			Name:           artistName,
			Albums:         albums[artistName].Slice(),
			RelatedArtists: relatedArtists,
		})
	}
	return artists, nil
}

func (db *ChordsDBv1Shim) Albums(_ context.Context, filters AlbumsFilters) ([]Album, error) {
	//TODO implement me
	panic("implement me")
}

func (db *ChordsDBv1Shim) Songs(_ context.Context, filters SongsFilters) ([]Song, error) {
	//TODO implement me
	panic("implement me")
}

var _ ChordsDBv1 = &ChordsDBv1Shim{}

// Convert an artist name into a "fake" ArtistID
func MakeArtistID(artistName string) ArtistID {
	return ArtistID(util.MakeID(artistName))
}

// Convert an album name into a "fake" AlbumID
func MakeAlbumID(albumName string) AlbumID {
	return AlbumID(util.MakeID(albumName))
}
