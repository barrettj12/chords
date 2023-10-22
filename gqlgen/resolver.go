package gqlgen

import (
	"github.com/barrettj12/chords/gqlgen/types"
	"github.com/barrettj12/chords/src/data"
)

// Resolver is the base class embedded inside all the GraphQL resolvers.
type Resolver struct {
	DB data.ChordsDBv1
}

// The following utility methods translate objects received from the data model
// into GraphQL API objects.

// translateArtist converts a data.Artist to a *types.Artist.
func (r *Resolver) translateArtist(artist data.Artist) *types.Artist {
	return &types.Artist{
		ID:   string(artist.ID),
		Name: artist.Name,
	}
}

// translateArtist converts a []data.Artist to a []*types.Artist.
func (r *Resolver) translateArtists(artistsData []data.Artist) []*types.Artist {
	artists := make([]*types.Artist, 0, len(artistsData))
	for _, artist := range artistsData {
		artists = append(artists, r.translateArtist(artist))
	}
	return artists
}

// translateAlbum converts a data.Album to a *types.Album.
func (r *Resolver) translateAlbum(album data.Album) *types.Album {
	return &types.Album{
		ID:   string(album.ID),
		Name: album.Name,
		Year: &album.Year,
	}
}

// translateAlbum converts a []data.Album to a []*types.Album.
func (r *Resolver) translateAlbums(albumsData []data.Album) []*types.Album {
	albums := make([]*types.Album, 0, len(albumsData))
	for _, album := range albumsData {
		albums = append(albums, r.translateAlbum(album))
	}
	return albums
}

// translateSong converts a data.Song to a *types.Song.
func (r *Resolver) translateSong(song data.Song) *types.Song {
	return &types.Song{
		ID:       string(song.ID),
		Name:     song.Name,
		TrackNum: &song.TrackNum,
		Chords:   string(song.Chords),
	}
}

// translateSong converts a []data.Song to a []*types.Song.
func (r *Resolver) translateSongs(songsData []data.Song) []*types.Song {
	songs := make([]*types.Song, 0, len(songsData))
	for _, song := range songsData {
		songs = append(songs, r.translateSong(song))
	}
	return songs
}
