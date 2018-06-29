package render

import (
	"net/http"
)

// Renderer ...
type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, status int, toRender, data interface{}) error
}
