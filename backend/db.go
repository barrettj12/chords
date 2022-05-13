// add license info here

// This file contains the database layer, which is responsible for
// communicating with the database, and converting between the database's
// format and Go objects.

package main

// getArtists retrieves the list of artists in the database.
func getArtists() []string {
	// TODO: implement
	return []string{"artist1", "artist2", "artist3"}
}

// album maps  song -> id
type album map[string]int

// songs maps albumname -> album
type songs map[string]album

// getSongs retrieves all songs by artist `artist`, grouped by album.
func getSongs(artist string) songs {
	// TODO: implement
	return songs{
		"album1": {
			"song1": 1,
			"song2": 2,
			"song3": 3,
		},
		"album2": {
			"song4": 4,
			"song5": 5,
		},
	}
}

// getChords retrieves the chords with the given id.
func getChords(id int) string {
	// TODO: implement
	return "Sample chords go here"
}
