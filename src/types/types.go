package types

type SongMeta struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Album    string `json:"album,omitempty"`
	TrackNum int    `json:"trackNum,omitempty"`
}
