package dblayer

import (
	"context"
	"log"

	gqltypes "github.com/barrettj12/chords/gqlgen/types"
)

// ChordsDBv1 is the data abstraction used for the chords v1 API.
type ChordsDBv1 interface {
	ArtistsV1(ctx context.Context) ([]*gqltypes.Artist, error)
}

func GetDBv1(url string, logger *log.Logger) ChordsDBv1 {
	logger.Printf("Using local filesystem database at %s\n", url)
	return NewLocalfs(url, logger)
}
