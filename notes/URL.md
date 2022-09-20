# Frontend URL structure

Although the backend API mostly uses song IDs to identify resources, e.g.
```
/api/v0/chords?id=1037
```
we want the frontend user interface to have nice, human-readable URLs. The artist and song name should be sufficient to uniquely identify any song.

### `mychords.com/`
The homepage will show all artists in the database, grouped alphabetically.

### `mychords.com/<artist>`
Show all songs by the given artist, grouped by album. For example,
```
mychords.com/elton-john
```

### `mychords.com/<artist>/<song>`
Show chords for the given song. For example,
```
mychords.com/elton-john/your-song
```
We will allow transposing the chords, and we can include a query parameter here:
```
mychords.com/elton-john/your-song?transpose=+3
```

## Notes
We are mapping artist/song names like this:
```
"Elton John" -> "elton-john"
"Your Song" -> "your-song"
```
We might need a table in the database to reverse this mapping.

We might also want to define abbreviations for unwieldy artist names. Compare
```
mychords.com/electric-light-orchestra
mychords.com/the-all-american-rejects
```
to
```
mychords.com/elo
mychords.com/aar
```