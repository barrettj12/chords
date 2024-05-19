package types

type SongMeta struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Album    string `json:"album,omitempty"`
	TrackNum int    `json:"trackNum,omitempty"`
}

type SearchResult struct {
	Type SearchResultType `json:"type"` // "song" or "artist"

	Name string    `json:"name,omitempty"` // populated for type:artist only
	ID   string    `json:"id,omitempty"`   // populated for type:song only
	Meta *SongMeta `json:"meta,omitempty"` // populated for type:song only
}

type SearchResultType string
