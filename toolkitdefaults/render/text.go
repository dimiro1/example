package render

import (
	"fmt"
	"net/http"
)

type Text struct{}

func (h Text) Render(w http.ResponseWriter, r *http.Request, status int, toRender interface{}) error {
	return h.RenderCtx(w, r, status, toRender, nil)
}

func (Text) RenderCtx(w http.ResponseWriter, r *http.Request, status int, toRender interface{}, context interface{}) error {
	w.Header().Set("Content-Type", "text/plain")
	var data []byte

	// Specific types
	switch toRender.(type) {
	case string:
		data = []byte(toRender.(string))
	case error:
		data = []byte(toRender.(error).Error())
	}

	// Stringer
	if s, ok := toRender.(fmt.Stringer); ok {
		data = []byte(s.String())
	}

	w.WriteHeader(status)
	_, err := w.Write(data)
	return err
}
