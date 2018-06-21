package render

import (
	"net/http"
)

// Renderer ...
type Renderer interface {
	// Render ...
	// context can be used to carry context data for HTML or text templates
	// note that this context is not from the context std library
	Render(w http.ResponseWriter, status int, toRender interface{}, context ...interface{}) error
}
