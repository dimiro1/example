package render

import (
	"net/http"
)

// Renderer ...
type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, status int, toRender interface{}) error
	RenderCtx(w http.ResponseWriter, r *http.Request, status int, toRender interface{}, context interface{}) error
}
