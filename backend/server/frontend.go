package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"

	"github.com/barrettj12/chords/backend/dblayer"
)

type Frontend struct {
	apiURL string
}

func (f *Frontend) artistsHandler(w http.ResponseWriter, r *http.Request) {
	resp, _ := http.Get(fmt.Sprintf("%s/api/v0/artists", f.apiURL))
	data, _ := io.ReadAll(resp.Body)
	artists := []string{}
	json.Unmarshal(data, &artists)
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
	resp, _ := http.Get(fmt.Sprintf("%s/api/v0/songs?artist=%s", f.apiURL, url.QueryEscape(artist)))
	data, _ := io.ReadAll(resp.Body)
	songs := []dblayer.SongMeta{}
	json.Unmarshal(data, &songs)

	w.Write([]byte(fmt.Sprintf("<h1>Songs by %s</h1>", artist)))
	for _, song := range songs {
		w.Write([]byte(
			fmt.Sprintf(`<a href="/b/chords?id=%[1]s">%[2]s</a><br>`, song.ID, song.Name),
		))
	}
}

func (f *Frontend) chordsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	metaResp, _ := http.Get(fmt.Sprintf("%s/api/v0/songs?id=%s", f.apiURL, url.QueryEscape(id)))
	data, _ := io.ReadAll(metaResp.Body)
	songs := []dblayer.SongMeta{}
	json.Unmarshal(data, &songs)
	if len(songs) == 0 {
		http.NotFound(w, r)
		return
	}

	chordsResp, _ := http.Get(fmt.Sprintf("%s/api/v0/chords?id=%s", f.apiURL, url.QueryEscape(id)))
	chords, _ := io.ReadAll(chordsResp.Body)

	w.Write([]byte(fmt.Sprintf("<h1>%s by %s</h1>", songs[0].Name, songs[0].Artist)))
	w.Write([]byte(fmt.Sprintf("<pre>%s</pre>", chords)))
}
