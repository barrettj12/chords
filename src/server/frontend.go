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

	"github.com/barrettj12/chords/src/client"
	"github.com/barrettj12/chords/src/dblayer"
	"github.com/barrettj12/chords/src/html"
	"github.com/barrettj12/chords/src/util"
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

func (f *Frontend) registerHandlers(mux http.ServeMux) {
	mux.HandleFunc("/b/artists", f.artistsHandler)
	mux.HandleFunc("/b/songs", f.songsHandler)
	mux.HandleFunc("/b/chords", f.chordsHandler)
	mux.HandleFunc("/b/random", f.randomHandler)

	// Default redirect to frontend artists page
	mux.Handle("/", http.RedirectHandler("/b/artists", http.StatusTemporaryRedirect))
}

func (f *Frontend) artistsHandler(w http.ResponseWriter, r *http.Request) {
	artists, _ := f.client.GetArtists()
	sortTitles(artists)

	body := html.Body{}
	body.Insert(html.NewHeading1("Artists"))

	p := html.NewParagraph()
	p.Insert(html.String("Click "))
	p.Insert(html.NewAnchor("/b/random", "here"))
	p.Insert(html.String(" for a random song."))
	body.Insert(p)

	// https://dev.to/jordanfinners/creating-a-collapsible-section-with-nothing-but-html-4ip9
	ul := html.NewUnorderedList()
	body.Insert(ul)
	for _, artist := range artists {
		ul.Insert(html.NewListItem(html.NewAnchor(
			fmt.Sprintf("/b/songs?artist=%s", url.QueryEscape(artist)),
			artist,
		)))
	}

	addFooter(&body)
	htmlDoc := html.HTML{
		Head: &html.Head{Title: "Jordy's Chordies"},
		Body: &body,
	}
	w.Write([]byte(htmlDoc.Render()))
}

// sortTitles sorts the given slice of titles, ignoring the articles "A", "An"
// and "The" at the beginning.
func sortTitles(titles []string) {
	sort.Slice(titles, func(i, j int) bool {
		return util.LessTitle(titles[i], titles[j])
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

	seeAlso, _ := f.client.SeeAlso(artist)
	if len(seeAlso) > 0 {
		body.Insert(html.NewHeading2("See also:"))

		ulSeeAlso := html.NewUnorderedList()
		body.Insert(ulSeeAlso)
		for _, relatedArtist := range seeAlso {
			ulSeeAlso.Insert(html.NewListItem(html.NewAnchor(
				fmt.Sprintf("/b/songs?artist=%s", url.QueryEscape(relatedArtist)),
				relatedArtist,
			)))
		}
	}

	addFooter(&body)
	htmlDoc := html.HTML{
		Head: &html.Head{Title: fmt.Sprintf("%s Chords", artist)},
		Body: &body,
	}
	w.Write([]byte(htmlDoc.Render()))
}

func (f *Frontend) chordsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	songs, _ := f.client.GetSongs(nil, &id, nil)
	if len(songs) == 0 {
		http.NotFound(w, r)
		return
	}
	chords, _ := f.client.GetChords(id)

	type pageData struct{ SongTitle, Artist, Chords, ID string }
	p := pageData{songs[0].Name, songs[0].Artist, string(chords), songs[0].ID}
	templateWithFooter := fmt.Sprintf(CHORDS_TEMPLATE, FOOTER)
	t := template.Must(template.New("page").Parse(templateWithFooter))
	t.Execute(w, p)
}

// TODO: the JavaScript code should go in a separate .js file
// Arguably, so should the HTML.
const CHORDS_TEMPLATE = `
<html>
  <head>
	 	<title>
		  {{.SongTitle}} Chords | {{.Artist}}
		</title>
	  <script type="module" src="https://barrettj12.github.io/chord-transposer/js/Main.js"></script>
	</head>
	<body>
	  <h1>
		  {{.SongTitle}} by {{.Artist}}
		</h1>
	  <pre id="chords" style="tab-size:3">loading...</pre>

		<b>Transpose:</b>
		<button id="minus">−</button>
    <input id="semitones" value="0" readonly="" style="text-align: center; width: 4.5ch;">
    <button id="plus">+</button>
		<button id="reset">Reset</button>

		<!-- footer goes here -->
		%s

		<script type="module">
			// Import JS backend
			import { transpose } from "https://barrettj12.github.io/chord-transposer/js/Main.js";

			// Get interactive elements on the page
			let chords = document.getElementById("chords");
			let plus = document.getElementById("plus");
			let minus = document.getElementById("minus");
			let semitones = document.getElementById("semitones");
			let reset = document.getElementById("reset");

			let originalChords = {{.Chords}};

			// Add event listeners
			reset.addEventListener("click", resetTranspose);
			plus.addEventListener("click", tuneUp);
			minus.addEventListener("click", tuneDown);

			// Check for cookie and restore previous transpose
			semitones.value = document.cookie.split("; ")
				.find((row) => row.startsWith("{{.ID}}="))?.split("=")[1] || 0;
			processChords();
		
			// Updates the textarea
			function processChords() {
				chords.innerHTML = transpose(originalChords, parseInt(semitones.value));
				document.cookie = "{{.ID}}="+semitones.value;
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

func (f *Frontend) randomHandler(w http.ResponseWriter, r *http.Request) {
	song, _ := f.client.RandomSong()
	http.Redirect(w, r, fmt.Sprintf("/b/chords?id=%s", url.QueryEscape(song.ID)), http.StatusSeeOther)
}

// addFooter adds a common footer to each page.
func addFooter(body *html.Body) {
	footer := html.String(FOOTER)
	body.Insert(footer)
}

const FOOTER = `
<footer>
	<hr>
	<p>© Jordan Barrett, 2023.</p>
	<p>Found a bug, problem or issue? Let me know <a href="https://github.com/barrettj12/chords/issues/new" target="_blank">here</a>.</p>
	<p>This project is open source - contribute <a href="https://github.com/barrettj12/chords" target="_blank">here</a>!</p>
</footer>`
