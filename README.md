# Jordy's chordies

This is (to be) a fully-functional web app showing chords that I've collected over the years.

The app (will) consist of several different parts:
- A frontend [web interface](https://barrettj12.github.io/chords/) hosted on [GitHub Pages](https://pages.github.com/). This is how end-users will access my chords.
- A backend web server, written in [Go](https://go.dev/) using the standard [`net/http`](https://pkg.go.dev/net/http) library.
- A database storing my chords in [an SQL format](notes/DATA_MODEL.md).
- A command-line interface, also written in Go, which I will privately use to update the chord database.

The app will interact with my own Chord Transposer API (coming soon) to provide transposition services.

The structure of this repo is as follows:
```bash
├─ backend/       # the Go backend
│  └─ dblayer/    # DB layer for backend
├─ cli/           # command-line tool
├─ docs/          # the frontend
└─ notes/         # Markdown specs/notes
```


## Motivation

I'm an amateur musician, and enjoy covering pop/rock songs. I've often found errors in chord sheets available online (including popular websites like [Ultimate Guitar](https://www.ultimate-guitar.com)). This led me to start working out chords myself, and over the years, I've amassed quite a collection.

Historically, I've stored these using a notes app like Google Keep. There are plenty of issues with this approach:
- Google Keep's only organisational structure is using "labels" - I can't sort/filter my chords by artist/album/etc.
- I often add [ASCII guitar tabs](https://en.wikipedia.org/wiki/ASCII_tab) to my chord sheets, to notate riffs. These really have to be rendered in a monospace font to ensure a nice layout and consistency. Of course, Google Keep doesn't use a monospace font.
- It's difficult to share my chords with others - I have to manually add them as a "collaborator" to each note.

Eventually, these concerns led me to decide that I'd be better off building my own website for my chords.