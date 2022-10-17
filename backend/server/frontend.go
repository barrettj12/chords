package server

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/barrettj12/chords/client"
)

type Frontend struct {
	apiURL string
	client *client.Client
}

func NewFrontend(apiURL string) (*Frontend, error) {
	c, err := client.NewClient(apiURL, "") // don't need authorisation
	if err != nil {
		return nil, err
	}

	return &Frontend{
		apiURL: apiURL,
		client: c,
	}, nil
}

func (f *Frontend) artistsHandler(w http.ResponseWriter, r *http.Request) {
	artists, _ := f.client.GetArtists()
	sort.Slice(artists, func(i, j int) bool { return artists[i] < artists[j] })

	w.Write([]byte("<h1>Artists</h1>"))
	for _, artist := range artists {
		w.Write([]byte(
			fmt.Sprintf(`<a href="/b/songs?artist=%[1]s">%[1]s</a><br>`, artist),
		))
	}
}

func (f *Frontend) songsHandler(w http.ResponseWriter, r *http.Request) {
	artist := r.URL.Query().Get("artist")
	songs, _ := f.client.GetSongs(&artist, nil, nil)

	w.Write([]byte(fmt.Sprintf("<h1>Songs by %s</h1>", artist)))
	for _, song := range songs {
		w.Write([]byte(
			fmt.Sprintf(`<a href="/b/chords?id=%[1]s">%[2]s</a><br>`, song.ID, song.Name),
		))
	}
}

func (f *Frontend) chordsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	songs, _ := f.client.GetSongs(nil, &id, nil)
	if len(songs) == 0 {
		http.NotFound(w, r)
		return
	}
	chords, _ := f.client.GetChords(id)

	w.Write([]byte(fmt.Sprintf("<h1>%s by %s</h1>", songs[0].Name, songs[0].Artist)))
	w.Write([]byte(fmt.Sprintf("<pre>%s</pre>", chords)))
}
