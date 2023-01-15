// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// src/server/frontend.go
// HTTP handlers for the frontend.

package server

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/barrettj12/chords/src/client"
	"github.com/barrettj12/chords/src/dblayer"
	"github.com/barrettj12/chords/src/html"
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
	sortTitles(artists)

	body := html.Body{}
	body.Insert(html.NewHeading1("Artists"))

	ul := html.NewUnorderedList()
	body.Insert(ul)
	for _, artist := range artists {
		ul.Insert(html.NewListItem(html.NewAnchor(
			fmt.Sprintf("/b/songs?artist=%s", url.QueryEscape(artist)),
			artist,
		)))
	}

	w.Write([]byte(body.Render()))
}

// sortTitles sorts the given slice of titles, ignoring the articles "A", "An"
// and "The" at the beginning.
func sortTitles(titles []string) {
	sort.Slice(titles, func(i, j int) bool {
		articles := []string{"A ", "An ", "The "}
		strip := func(s string) string {
			for _, a := range articles {
				s = strings.TrimPrefix(s, a)
			}
			return s
		}

		return strip(titles[i]) < strip(titles[j])
	})
}

func (f *Frontend) songsHandler(w http.ResponseWriter, r *http.Request) {
	artist := r.URL.Query().Get("artist")
	songs, _ := f.client.GetSongs(&artist, nil, nil)

	// Group songs by album
	albums := map[string][]dblayer.SongMeta{}
	for _, song := range songs {
		albums[song.Album] = append(albums[song.Album], song)
		sort.Slice(albums[song.Album], func(i, j int) bool {
			return albums[song.Album][i].TrackNum < albums[song.Album][j].TrackNum
		})
	}
	// Sort albums
	albumNames := []string{}
	for album := range albums {
		albumNames = append(albumNames, album)
	}
	sort.Slice(albumNames, func(i, j int) bool {
		if albumNames[j] == "" && albumNames[i] != "" {
			return true
		}
		if albumNames[i] == "" && albumNames[j] != "" {
			return false
		}
		return albumNames[i] < albumNames[j]
	})

	// Construct and render HTML
	body := html.Body{}
	body.Insert(html.NewHeading1(fmt.Sprintf("Songs by %s", artist)))

	ulAlbums := html.NewUnorderedList()
	body.Insert(ulAlbums)
	for _, album := range albumNames {
		tracklist := albums[album]
		var htmlListTracks html.List
		if album == "" {
			album = "(no album)"
			htmlListTracks = html.NewUnorderedList()
		} else {
			htmlListTracks = html.NewOrderedList()
		}
		ulAlbums.Insert(html.NewListItem(html.String(album), htmlListTracks))

		for _, song := range tracklist {
			li := html.NewListItem(html.NewAnchor(
				fmt.Sprintf("/b/chords?id=%s", url.QueryEscape(song.ID)),
				song.Name,
			))
			li.SetValue(fmt.Sprint(song.TrackNum))
			htmlListTracks.Insert(li)
		}
	}

	w.Write([]byte(body.Render()))
}

func (f *Frontend) chordsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	songs, _ := f.client.GetSongs(nil, &id, nil)
	if len(songs) == 0 {
		http.NotFound(w, r)
		return
	}
	chords, _ := f.client.GetChords(id)

	type pageData struct{ SongTitle, Artist, Chords string }
	p := pageData{songs[0].Name, songs[0].Artist, string(chords)}
	t := template.Must(template.New("page").Parse(CHORDS_TEMPLATE))
	t.Execute(w, p)
}

const CHORDS_TEMPLATE = `
<html>
  <head>
	  <script type="module" src="https://barrettj12.github.io/chord-transposer/js/Main.js"></script>
	</head>
	<body>
	  <h1>
		  {{.SongTitle}} by {{.Artist}}
		</h1>
	  <pre id="chords">{{.Chords}}</pre>

		<b>Transpose:<b>
		<button id="minus">−</button>
    <input id="semitones" value="0" readonly="" style="text-align: center; width: 4.5ch;">
    <button id="plus">+</button>
		<button id="reset">Reset</button>

		<script type="module">
        // Import JS backend
        import { transpose } from "https://barrettj12.github.io/chord-transposer/js/Main.js";

        // Get interactive elements on the page
        let chords = document.getElementById("chords");
        let plus = document.getElementById("plus");
        let minus = document.getElementById("minus");
        let semitones = document.getElementById("semitones");
        let reset = document.getElementById("reset");

				let originalChords = chords.innerHTML;

        // Add event listeners
        reset.addEventListener("click", resetTranspose);
        plus.addEventListener("click", tuneUp);
        minus.addEventListener("click", tuneDown);
      
        // Updates the textarea
        function processChords() {
          chords.innerHTML = transpose(originalChords, parseInt(semitones.value));
        }

        function resetTranspose() {
          semitones.value = 0;
          processChords();
        }

        function tuneUp() {
          semitones.value = (parseInt(semitones.value) + 1).toString();
          processChords();
        }

        function tuneDown() {
          semitones.value = (parseInt(semitones.value) - 1).toString();
          processChords();
        }
      </script> 
  </body>
</html>
`
