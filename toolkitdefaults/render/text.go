package render

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Text struct{}

func (Text) Render(w http.ResponseWriter, r *http.Request, status int, toRender interface{}, _ interface{}) error {
	if w == nil {
		return errors.New("render: http.ResponseWriter cannot be nil")
	}

	if r == nil {
		return errors.New("render: *http.Request cannot be nil")
	}

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
	return errors.WithStack(err)
}
