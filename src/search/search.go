package search

import (
	"strings"

	"github.com/barrettj12/chords/src/types"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
)

type Index struct {
	bleveIndex bleve.Index
}

// NewIndex returns a new, empty Index.
func NewIndex() (*Index, error) {
	mapping := bleve.NewIndexMapping()
	// TODO: use bleve.New instead and persist the index on disk
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return nil, err
	}

	return &Index{
		bleveIndex: index,
	}, nil
}

func (i *Index) Add(meta types.SongMeta) error {
	return i.bleveIndex.Index(meta.ID, meta)
}

func (i *Index) Remove(id string) error {
	return i.bleveIndex.Delete(id)
}

func (i *Index) Search(rawQuery string) (ids []string, err error) {
	// For some reason, terms are not matched with mixed case
	// So map everything to lowercase
	rawQuery = strings.ToLower(rawQuery)

	words := strings.Split(rawQuery, " ")
	termQueries := make([]query.Query, 0, len(words))
	for _, w := range words {
		if w == "" {
			continue
		}
		termQueries = append(termQueries, bleve.NewPrefixQuery(w))
	}

	query := bleve.NewConjunctionQuery(termQueries...)
	search := bleve.NewSearchRequest(query)

	searchResults, err := i.bleveIndex.Search(search)
	if err != nil {
		return nil, err
	}

	for _, res := range searchResults.Hits {
		ids = append(ids, res.ID)
	}
	return ids, nil
}
