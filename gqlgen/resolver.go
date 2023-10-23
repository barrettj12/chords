package gqlgen

import (
	"github.com/barrettj12/chords/gqlgen/types"
	"github.com/barrettj12/chords/src/data"
)

// Resolver is the base class embedded inside all the GraphQL resolvers.
type Resolver struct {
	DB data.ChordsDBv1
}

// General purpose resolvers, which

// resolveArtist converts a []data.Artist (which should only have 0 or 1
// elements) into a *types.Artist, and also handles errors from the DB.
func (r *Resolver) resolveArtist(artistsData []data.Artist, err error) (*types.Artist, error) {
	if err != nil {
		return nil, err
	}

	if len(artistsData) == 0 {
		return nil, nil
	}
	return r.translateArtist(artistsData[0]), nil
}

// resolveArtists converts a []data.Artist into a []*types.Artist, and also
// handles errors from the DB.
func (r *Resolver) resolveArtists(artistsData []data.Artist, err error) ([]*types.Artist, error) {
	if err != nil {
		return nil, err
	}

	artists := make([]*types.Artist, 0, len(artistsData))
	for _, artist := range artistsData {
		artists = append(artists, r.translateArtist(artist))
	}
	return artists, nil
}

// resolveAlbum converts a []data.Album (which should only have 0 or 1
// elements) into a *types.Album, and also handles errors from the DB.
func (r *Resolver) resolveAlbum(albumsData []data.Album, err error) (*types.Album, error) {
	if err != nil {
		return nil, err
	}

	if len(albumsData) == 0 {
		return nil, nil
	}
	return r.translateAlbum(albumsData[0]), nil
}

// resolveAlbums converts a []data.Album into a []*types.Album, and also
// handles errors from the DB.
func (r *Resolver) resolveAlbums(albumsData []data.Album, err error) ([]*types.Album, error) {
	if err != nil {
		return nil, err
	}

	albums := make([]*types.Album, 0, len(albumsData))
	for _, album := range albumsData {
		albums = append(albums, r.translateAlbum(album))
	}
	return albums, nil
}

// resolveSong converts a []data.Song (which should only have 0 or 1
// elements) into a *types.Song, and also handles errors from the DB.
func (r *Resolver) resolveSong(songsData []data.Song, err error) (*types.Song, error) {
	if err != nil {
		return nil, err
	}

	if len(songsData) == 0 {
		return nil, nil
	}
	return r.translateSong(songsData[0]), nil
}

// resolveSongs converts a []data.Song into a []*types.Song, and also handles
// errors from the DB.
func (r *Resolver) resolveSongs(songsData []data.Song, err error) ([]*types.Song, error) {
	if err != nil {
		return nil, err
	}

	songs := make([]*types.Song, 0, len(songsData))
	for _, song := range songsData {
		songs = append(songs, r.translateSong(song))
	}
	return songs, nil
}

// Stateless methods which simply translate the types from the data package
// into corresponding GraphQL API types.

// translateArtist converts a data.Artist into a *types.Artist.
func (r *Resolver) translateArtist(artist data.Artist) *types.Artist {
	return &types.Artist{
		ID:   string(artist.ID),
		Name: artist.Name,
	}
}

// translateAlbum converts a data.Album into a *types.Album.
func (r *Resolver) translateAlbum(album data.Album) *types.Album {
	return &types.Album{
		ID:   string(album.ID),
		Name: album.Name,
		Year: &album.Year,
	}
}

// translateSong converts a data.Song into a *types.Song.
func (r *Resolver) translateSong(song data.Song) *types.Song {
	return &types.Song{
		ID:       string(song.ID),
		Name:     song.Name,
		TrackNum: &song.TrackNum,
		Chords:   string(song.Chords),
	}
}
