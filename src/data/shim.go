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

// TODO: handle RelatedTo filter
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
	// Get album info from song metadata
	songs, err := db.db.GetSongs("", "", "")
	if err != nil {
		return nil, err
	}

	// map album name -> Album struct
	albums := map[string]*Album{}
	for _, song := range songs {
		if song.Album == "" {
			continue
		}

		if albums[song.Album] == nil {
			// Make new Album
			albums[song.Album] = &Album{
				ID:     MakeAlbumID(song.Album),
				Name:   song.Album,
				Artist: MakeArtistID(song.Artist),
				Songs:  []SongID{},
			}
		}

		albums[song.Album].Songs = append(albums[song.Album].Songs, SongID(song.ID))
	}

	albumsSlice := []Album{}
	for _, album := range albums {
		albumsSlice = append(albumsSlice, *album)
	}
	return albumsSlice, nil
}

func (db *ChordsDBv1Shim) Songs(_ context.Context, filters SongsFilters) ([]Song, error) {
	rawSongs, err := db.db.GetSongs("", "", "")
	if err != nil {
		return nil, err
	}

	songs := []Song{}
	for _, song := range rawSongs {
		chords, _ := db.db.GetChords(song.ID)

		songs = append(songs, Song{
			ID:       SongID(song.ID),
			Name:     song.Name,
			Artist:   MakeArtistID(song.Artist),
			Album:    MakeAlbumID(song.Album),
			TrackNum: song.TrackNum,
			Chords:   chords,
		})
	}
	return songs, nil
}

// Convert an artist name into a "fake" ArtistID
func MakeArtistID(artistName string) ArtistID {
	return ArtistID(util.MakeID(artistName))
}

// Convert an album name into a "fake" AlbumID
func MakeAlbumID(albumName string) AlbumID {
	return AlbumID(util.MakeID(albumName))
}
