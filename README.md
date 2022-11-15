# Jordy's chordies

This is a web app showing chords that I've collected / worked out over the
years. You can find the live app at https://chords.fly.dev/.


## In a bit more detail...

The app consists of several pieces.
- The centrepiece is an API server, which provides an [API](docs/API.md) to
  get, add, update and delete chords.
- A [frontend](https://chords.fly.dev/) giving a nice interface to find and
  view the chords. This also includes the ability to transpose chords (reusing
  code from my [chord-transposer](https://github.com/barrettj12/chord-transposer)
  project).
- A persistent file system which the API server uses to store the chords and
  their metadata.

All of the above are hosted on [Fly.io](https://fly.io/). The API server and
frontend server are both written in [Go](https://go.dev/) using the standard
[`net/http`](https://pkg.go.dev/net/http) library.

I've also written a command-line interface in Go, which I privately use to add
and update chords.


## The structure of this repository

- `docs`: Markdown specs and explanatory notes. Read these if you'd like to
  learn more about the inner workings of the app.
- `src`: source code
  - `client`: API client, used by the CLI and frontend
  - `cmd`: CLI - I use this to update the chords database
  - `dblayer`: core data structures and database wrappers, used by the client
    and server
  - `server`: the API and frontend server
- `tests`: Go integration tests for the API server.
- `main.go`: main entry point, runs the API/frontend server.


## Motivation

I'm an amateur musician, and enjoy covering pop/rock songs. I've often found
errors in chord sheets available online (including popular websites like
[Ultimate Guitar](https://www.ultimate-guitar.com)). This led me to start
working out chords myself, and over the years, I've amassed quite a collection.

Historically, I've stored these using a notes app like Google Keep. There are
plenty of issues with this approach:
- Google Keep's only organisational structure is using "labels" - I can't
sort/filter my chords by artist/album/etc.
- I often add [ASCII guitar tabs](https://en.wikipedia.org/wiki/ASCII_tab) to
my chord sheets, to notate riffs. These really have to be rendered in a
monospace font to ensure a nice layout and consistency. Of course, Google Keep
doesn't use a monospace font.
- It's difficult to share my chords with others - I have to manually add them
as a "collaborator" to each note.

Eventually, these concerns led me to decide that I'd be better off building my
own website for my chords.


## License

This project is open-source, and licensed under the terms of the
[GNU Affero General Public License](https://www.gnu.org/licenses/agpl-3.0.en.html).