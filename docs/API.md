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


### `GET /view?id=<id>`

Returns chords for the song with id `<id>`, as plain text.

### `POST /update?id=<id>`
*Should we use the POST or PUT method?*

Sets the chords with id `<id>` to the value of the request body. **This method requires authentication**.

*Maybe we should also support getting/setting by specifying artist, album, song name, e.g.*
```
GET /view?artist=<artist>&song=<song>
PUT /update?artist=<artist>&song=<song>
```

### `GET /search?q=<query>`

Searches the database for songs matching `<query>`, and returns a list of matching songs.