package data

type Artist struct {
	ID             ArtistID   `json:"id"`
	Name           string     `json:"name"`
	Albums         []AlbumID  `json:"albums"`
	RelatedArtists []ArtistID `json:"relatedArtists"`
}

type ArtistID string

type Album struct {
	ID     AlbumID  `json:"id"`
	Name   string   `json:"name"`
	Year   int      `json:"year,omitempty"`
	Artist ArtistID `json:"artist"`
	Songs  []SongID `json:"songs"`
}

type AlbumID string

type Song struct {
	ID       SongID   `json:"id"`
	Name     string   `json:"name"`
	Artist   ArtistID `json:"artist,omitempty"`
	Album    AlbumID  `json:"album,omitempty"`
	TrackNum int      `json:"trackNum,omitempty"`
	Chords   []byte   `json:"chords"`
}

type SongID string
