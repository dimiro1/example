package app

import (
	"net/http"
)

// index render the root page
func (a *Application) index() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.renderer.Render(w, r, http.StatusOK, "Welcome")
	})
}
