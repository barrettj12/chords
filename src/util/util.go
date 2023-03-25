// Jordy's Chordies - a web app for song chords
//     https://github.com/barrettj12/chords
// Copyright 2022, Jordan Barrett (@barrettj12)
//     https://github.com/barrettj12
// Licensed under the GNU AGPLv3.

// src/util/util.go
// Utility functions used in other parts of the code.

package util

import "strings"

// LessTitle compares two titles and returns true if title1 is alphabetically
// before title2, ignoring preceding articles (e.g. "a", "an", "the").
// It can be used to sort a slice as follows:
//
//	sort.Slice(titles, func(i, j int) bool {
//		return util.LessTitle(titles[i], titles[j])
//	})
func LessTitle(title1, title2 string) bool {
	articles := []string{"A ", "An ", "The "}
	strip := func(s string) string {
		for _, a := range articles {
			s = strings.TrimPrefix(s, a)
		}
		return s
	}

	return strip(title1) < strip(title2)
}
