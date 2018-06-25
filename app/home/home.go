package home

import (
	"net/http"

	"github.com/dimiro1/example/toolkit/render"
	"github.com/dimiro1/example/toolkit/router"
)

type Home struct {
	renderer render.Renderer
}

func (h *Home) Name() string {
	return "home"
}

func (h *Home) RegisterRoutes(router router.Router) {
	router.HandleFunc("GET", "/", h.index())
}

func NewHome(render render.Renderer) *Home {
	return &Home{render}
}

// index render the root page
func (h *Home) index() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.renderer.Render(w, http.StatusOK, "Welcome")
	})
}
