// add license info here

// add description of file here

package main

import (
	"log"
	"net/http"
)

func main() {
	// Register API endpoints
	http.HandleFunc("/list", listHandler)     // list available chord sheets
	http.HandleFunc("/view", viewHandler)     // view a given chord sheet
	http.HandleFunc("/update", updateHandler) // update a chord sheet

	// Start listening on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handles requests to list available chord sheets.
func listHandler(w http.ResponseWriter, r *http.Request) {
	// parse request `r`
	// send to list function
	// write output to `w`
}

// Handles requests to view a given chord sheet.
func viewHandler(w http.ResponseWriter, r *http.Request) {
	// parse request `r`
	// send to get function
	// write output to `w`
}

// Handles requests to update a given chord sheet.
func updateHandler(w http.ResponseWriter, r *http.Request) {
	// parse request `r`
	// send to update function
	// write output to `w`
}
