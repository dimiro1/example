package home

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/dimiro1/example/log"
	"github.com/dimiro1/example/toolkit/render"
	"github.com/dimiro1/example/toolkit/router"
)

type Home struct {
	logger   *log.Logger
	renderer render.Renderer
}

func (h *Home) Name() string {
	return "home"
}

func (h *Home) RegisterRoutes(router router.Router) {
	router.HandleFunc("GET", "/", h.index())
}

func NewHome(logger *log.Logger, render render.Renderer) (*Home, error) {
	if logger == nil {
		return nil, errors.New("recipes: logger cannot be nil")
	}

	if render == nil {
		return nil, errors.New("recipes: render cannot be nil")
	}

	return &Home{logger: logger, renderer: render}, nil
}

// index render the root page
func (h *Home) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.renderer.Render(w, r, http.StatusOK, "index.tmpl", "Welcome to example"); err != nil {
			h.logger.ErrorRendering(err, "Home.index")
		}
	}
}
