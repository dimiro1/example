package render

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/pkg/errors"
)

type HTML struct {
	templates *template.Template
}

func (h *HTML) Render(w http.ResponseWriter, r *http.Request, status int, toRender interface{}) error {
	return h.RenderCtx(w, r, status, toRender, nil)
}

func (h *HTML) RenderCtx(w http.ResponseWriter, r *http.Request, status int, toRender interface{}, context interface{}) error {
	if toRender == nil {
		return errors.New("render: toRender cannot be nil")
	}

	if _, ok := toRender.(string); !ok {
		return errors.New("render: toRender must be a string")
	}

	if t := h.templates.Lookup(toRender.(string)); t == nil {
		return fmt.Errorf("render: the template %s could not be found", toRender.(string))
	}

	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	w.WriteHeader(status)
	return h.templates.ExecuteTemplate(w, toRender.(string), context)
}

func NewHTML(templates *template.Template) *HTML {
	return &HTML{
		templates: templates,
	}
}
