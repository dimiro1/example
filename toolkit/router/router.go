package router

import (
	"net/http"
)

// Router basic http router definition
type Router interface {
	http.Handler

	// Handle registers a new route with a matcher for the URL path and method.
	Handle(method, path string, handler http.Handler)

	// HandleFunc registers a new route with a matcher for the URL path and method.
	HandleFunc(method, path string, handler http.HandlerFunc)

	// NotFound handler to be used when there are no routes
	NotFound(handler http.Handler)
}
