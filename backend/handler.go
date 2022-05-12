// add license here

// Defines a type Handler which

package handler

import "net/http"

type Handler struct {
	methods []string
}

// unnecessary
// func (h *Handler) addRoute(url string) {
// 	http.Handle(url, h)
// }

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle OPTIONS request
	if r.Method == http.MethodOptions {
		// Add allowed methods to Allow header
		for _, method := range h.methods {
			w.Header.Add("Allow", method)
		}
		return
	}

	// Check method is allowed
	if !h.allows(r.Method) {

	}
}
