package server

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"github.com/barrettj12/chords/src/client"
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
