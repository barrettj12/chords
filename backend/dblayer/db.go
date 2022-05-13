// add license info here

// This file contains the database layer, which is responsible for
// communicating with the database, and converting between the database's
// format and Go objects.

package dblayer

// getArtists retrieves the list of artists in the database.
func getArtists() []string {
	// TODO: implement
	return []string{"artist1", "artist2", "artist3"}
}

// getSongs retrieves all songs by artist `artist`, grouped by album.
func getSongs(artist string) Songs {
	// TODO: implement
	return Songs{
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
