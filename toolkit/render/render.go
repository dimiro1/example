package render

import (
	"net/http"
)

// Renderer ...
type Renderer interface {
	// Render ...
	// extra can be used to carry context data for HTML or text templates
	Render(w http.ResponseWriter, status int, data interface{}, extra ...interface{}) error
}
