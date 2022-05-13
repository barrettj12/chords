// add license here

// Defines a type Handler which

// putting unnecessary effort in here

package main

import "net/http"

// List all HTTP methods
var METHODS = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

type Handler struct {
	handleGet func()
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
		// for _, method := range h.methods {
		// 	w.Header.Add("Allow", method)
		// }
		// return
	}

	// Catch unallowed methods
	if !h.allows(r.Method) {
		// http.
	}

	//	method string -> methodHandler func(...) ...
}

// handleMethod returns the handler function for the provided method.
func (h *Handler) handleMethod(method string) interface{} {
	switch method {
	case http.MethodGet:
		return h.handleGet
	default:
		return nil
	}
}

// methods gets the methods which are allowed for this handler.
func (h *Handler) methods() []string {
	methods := []string{}
	for _, method := range METHODS {
		if h.allows(method) {
			methods = append(methods, method)
		}
	}
	return methods
}

// allows
func (h *Handler) allows(method string) bool {
	return h.handleMethod(method) != nil
}
