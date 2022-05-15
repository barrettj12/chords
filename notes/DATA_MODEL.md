# The database model

All the data will be stored in a single table `chords`, with the following
columns:

| id  | artist | album | song | data |
|-----|--------|-------|------|------|
| NUM |  CHAR  | CHAR  | CHAR | TEXT |

*Do we want to add extra metadata, e.g. year?*

We might consider adding additional tables/data structures to make other queries efficient, e.g. browsing by artist, searching. In this case, the `chords` table will be considered the "one true source" of all data, and the other tables/data structures will be derived from it.

If we do this, we'll also add a `validate()` function in `db.go` that reinitialises all the derived tables from the `chords` table. This will be run periodically to ensure consistency of the data (with option to run manually).