// add license info here

// add description of file here

package main

import (
	"github.com/barrettj12/jordys-chordies/backend/dblayer"
	"github.com/barrettj12/jordys-chordies/backend/server"
)

func main() {
	db := &dblayer.TempDB{}
	s := server.New(db, ":8080")
	s.Start()
}
