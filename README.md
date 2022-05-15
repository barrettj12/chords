# Jordy's chordies

This is (to be) a fully-functional web app showing chords that I've collected over the years.

As an amateur musician, I enjoy figuring out the chords to a song. I often find errors in chord sheets available online (including popular websites like [Ultimate Guitar](https://www.ultimate-guitar.com)).  

The app (will) consist of several different parts:
- A frontend web interface hosted on GitHub Pages. This is how end-users will access my chords.
- A backend web server, written in Go using the standard `net/http` library.
- A database storing my chords in an SQL format.
- A command-line interface, which I will privately use to update the chord database.

The app will interact with my own Chord Transposer API (coming soon) to provide transposition services.