package render

import (
	"html/template"
	"net/http"

	"github.com/pkg/errors"
)

type HTML struct {
	templates *template.Template
}

func (h *HTML) Render(w http.ResponseWriter, r *http.Request, status int, toRender interface{}, data interface{}) error {
	if w == nil {
		return errors.New("render: http.ResponseWriter cannot be nil")
	}

	if r == nil {
		return errors.New("render: *http.Request cannot be nil")
	}

	if toRender == nil {
		return errors.New("render: toRender cannot be nil")
	}

	if _, ok := toRender.(string); !ok {
		return errors.New("render: toRender must be a string")
	}

	if t := h.templates.Lookup(toRender.(string)); t == nil {
		return errors.Errorf("render: the template %s could not be found", toRender.(string))
	}

	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(status)
	return errors.WithStack(h.templates.ExecuteTemplate(w, toRender.(string), data))
}

func NewHTML(templates *template.Template) *HTML {
	return &HTML{
		templates: templates,
	}
}
