package home

import (
	"errors"
	"net/http"

	"github.com/dimiro1/example/toolkit/render"
	"github.com/dimiro1/example/toolkit/router"
	log "github.com/sirupsen/logrus"
)

type Home struct {
	logger   *log.Entry
	renderer render.Renderer
}

func (h *Home) Name() string {
	return "home"
}

func (h *Home) RegisterRoutes(router router.Router) {
	router.HandleFunc("GET", "/", h.index())
}

func NewHome(logger *log.Entry, render render.Renderer) (*Home, error) {
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := h.renderer.RenderCtx(w, r, http.StatusOK, "index.tmpl", "Welcome to example"); err != nil {
			h.logger.WithError(err).Error("searchRecipes")
		}
	})
}
