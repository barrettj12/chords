// add license info here

// This file contains the database layer, which is responsible for
// communicating with the database, and converting between the database's
// format and Go objects.

package main

// getArtists retrieves the list of artists in the database.
func getArtists() []string {
	// TODO: implement
	return nil
}

// album maps  song -> id
type album map[string]int

// songs maps albumname -> album
type songs map[string]album

// getSongs retrieves all songs by artist `artist`, grouped by album.
func getSongs(artist string) songs {
	// TODO: implement
	return nil
}

// getChords retrieves the chords with the given id.
func getChords(id int) string {
	return ""
}
