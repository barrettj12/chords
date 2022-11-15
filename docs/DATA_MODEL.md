# The database model

I am using a file tree to store the chord data. The file structure is
as follows:
```
basedir
├─ [id1]
│  ├─ meta.json
│  └─ chords.txt
├─ [id2]
│  ├─ meta.json
│  └─ chords.txt
...
```

In words: every song in the DB has a unique ID. The ID is used as the name of
a corresponding subfolder inside the root directory. Inside this subfolder,
there are two files: `meta.json` and `chords.txt`

The `meta.json` file is a JSON file containing the metadata for the song
(song name, artist, etc). The format of this file is like:
```json
{
  "id": "BananaPancakes",
  "name": "Banana Pancakes",
  "artist": "Jack Johnson",
  "album": "In Between Dreams",
  "trackNum": 3
}
```

The `"id"` must be identical to the name of the subfolder. `"trackNum"` is
the position in which the song appears on its album - this is used to display
albums correctly on the frontend. The other fields are self-explanatory.

`chords.txt` simply contains the chords in plain-text format.


## Alternative relational model

*I decided against using SQL for simplicity, but I'll keep this section around
for interest's sake, and also in case I decide to switch to an SQL DB in future.*

All the data will be stored in a single table `chords`, with the following
columns:

| id  | artist | album | song | data |
|-----|--------|-------|------|------|
| NUM |  CHAR  | CHAR  | CHAR | TEXT |

*Do we want to add extra metadata, e.g. year?*

We might consider adding additional tables/data structures to make other queries
efficient, e.g. browsing by artist, searching. In this case, the `chords` table
will be considered the "one true source" of all data, and the other tables/data
structures will be derived from it.

If we do this, we'll also add a `validate()` function in `db.go` that
reinitialises all the derived tables from the `chords` table. This will be run
periodically to ensure consistency of the data (with option to run manually).