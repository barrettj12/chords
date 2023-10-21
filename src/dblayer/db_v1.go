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

func GetDBv1(url string, logger *log.Logger) (ChordsDBv1, error) {
	db, err := GetDB(url, logger)
	if err != nil {
		return nil, err
	}
	return &ChordsDBv1Shim{db}, nil
}

// ChordsDBv1Shim is an adapter from a ChordsDB to a ChordsDBv1.
type ChordsDBv1Shim struct {
	db ChordsDB
}

func (db *ChordsDBv1Shim) ArtistsV1(ctx context.Context) ([]*gqltypes.Artist, error) {
	artistNames, err := db.db.GetArtists()
	if err != nil {
		return nil, err
	}

	//// Generate see also map
	//seeAlsoMap := map[string][]string{}
	//
	//seeAlsoPath := filepath.Join(l.basedir, "see-also.json")
	//seeAlsoFile, err := os.Open(seeAlsoPath)
	//seeAlsos := [][]string{}
	//err = json.NewDecoder(seeAlsoFile).Decode(&seeAlsos)
	//if err != nil {
	//	l.log.Printf("WARNING couldn't unmarshal see also data: %v", err)
	//}
	//
	//for _, grp := range seeAlsos {
	//	seeAlsoMap[grp[0]] = append(seeAlsoMap[grp[0]], grp[1])
	//	seeAlsoMap[grp[1]] = append(seeAlsoMap[grp[1]], grp[0])
	//}

	// Generate returned data
	var artists []*gqltypes.Artist
	for _, name := range artistNames {
		artists = append(artists, &gqltypes.Artist{
			Name: name,
			// TODO: how to put in albums / related artists?
		})
	}
	return artists, nil
}
