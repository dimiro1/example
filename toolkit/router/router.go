package router

import (
	"net/http"
)

// Route represents a registered Route
type Route struct {
	// Method HTTP Method, GET, POST, DELETE ...
	Method string
	// Path rhe registered URL Path
	Path string
	// Handler the actual handler
	Handler http.Handler
	// HandlerName holds the path/name of the handler inside your code
	HandlerName string
}

// Router basic http router definition
type Router interface {
	http.Handler

	// Handle registers a new route with a matcher for the URL path and method.
	Handle(method, path string, handler http.Handler)

	// HandleFunc registers a new route with a matcher for the URL path and method.
	HandleFunc(method, path string, handler http.HandlerFunc)

	// NotFound handler to be used when there are no routes
	NotFound(handler http.Handler)

	// Routes returns a list of registered routes
	Routes() []Route
}
