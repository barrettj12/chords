package gqlgen

import "github.com/barrettj12/chords/src/dblayer"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB dblayer.ChordsDBv1
}
