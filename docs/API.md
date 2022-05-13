# The backend API interface


### `GET /artists`

Returns a JSON list of all artists in the database.


### `GET /songs?artist=<artist>`

Returns a JSON object containing all songs by `<artist>`, grouped by album:
```
{
  "album1": {
    "song1": <id>,
    "song2": <id>,
    ...
  },
  "album2": {
    ...
  },
  ...
}
```


### `GET /chords/<id>`

Returns chords for the song with id `<id>`, as plain text.


### `PUT /chords/<id>`

Sets the chords with id `<id>` to the value of the request body. **This method requires authentication**.


### `POST /chords`

Adds a new chord sheet to the database. **This method requires authentication**.

The request body should be a JSON object of the following form:
```
{
  "artist": ...
  "album":  ...
  "song":   ...
  "chords": ...
}
```


### `GET /search?q=<query>`

Searches the database for songs matching `<query>`, and returns a JSON object containing matching songs and ids.