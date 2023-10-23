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
	songMap := map[string]*collections.Set[SongID]{}
	for _, song := range songs {
		if song.Artist == "" {
			continue
		}

		if albums[song.Artist] == nil {
			albums[song.Artist] = collections.NewSet[AlbumID](0)
		}
		albums[song.Artist].Add(MakeAlbumID(song.Album))

		if songMap[song.Artist] == nil {
			songMap[song.Artist] = collections.NewSet[SongID](0)
		}
		songMap[song.Artist].Add(SongID(song.ID))
	}

	artists := []Artist{}
	for _, artistName := range artistNames {
		id := MakeArtistID(artistName)
		if filters.ID != "" && id != filters.ID {
			// We requested a specific artist by ID, this is not it
			continue
		}
		if filters.Album != "" && !albums[artistName].Contains(filters.Album) {
			// We requested the artist who created a given album, this is not it
			continue
		}
		if filters.Song != "" && !songMap[artistName].Contains(filters.Song) {
			// We requested the artist who created a given song, this is not it
			continue
		}

		// Get "see also" data to fill related artists
		seeAlso, _ := db.db.SeeAlso(artistName)

		relatedArtists := []ArtistID{}
		relatedToFilterPassed := false
		for _, relatedArtist := range seeAlso {
			relatedID := MakeArtistID(relatedArtist)
			relatedArtists = append(relatedArtists, relatedID)

			if relatedID == filters.RelatedTo {
				relatedToFilterPassed = true
			}
		}
		if filters.RelatedTo != "" && !relatedToFilterPassed {
			// We've been asked for artists related to X, but this artist is not
			// related to X. So skip it.
			continue
		}

		artists = append(artists, Artist{
			ID:             id,
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
		id := MakeAlbumID(song.Album)
		if filters.ID != "" && id != filters.ID {
			continue
		}

		artistID := MakeArtistID(song.Artist)
		if filters.Artist != "" && artistID != filters.Artist {
			// We requested albums for a given artist, which doesn't match this album
			continue
		}

		if albums[song.Album] == nil {
			// Make new Album
			albums[song.Album] = &Album{
				ID:     id,
				Name:   song.Album,
				Artist: artistID,
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
	rawSongs, err := db.db.GetSongs("", string(filters.ID), "")
	if err != nil {
		return nil, err
	}

	songs := []Song{}
	for _, song := range rawSongs {
		albumID := MakeAlbumID(song.Album)
		if filters.Album != "" && albumID != filters.Album {
			// We requested songs for a given album, which doesn't include this song
			continue
		}

		chords, _ := db.db.GetChords(song.ID)

		songs = append(songs, Song{
			ID:       SongID(song.ID),
			Name:     song.Name,
			Artist:   MakeArtistID(song.Artist),
			Album:    albumID,
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
