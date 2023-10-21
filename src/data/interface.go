package data

import "context"

type ChordsDBv1 interface {
	Artists(context.Context, ArtistsFilters) ([]Artist, error)
	Albums(context.Context, AlbumsFilters) ([]Album, error)
	Songs(context.Context, SongsFilters) ([]Song, error)
}

type ArtistsFilters struct{}

type AlbumsFilters struct{}

type SongsFilters struct{}
