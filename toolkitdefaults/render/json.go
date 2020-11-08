package render

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type JSON struct{}

func (JSON) Render(w http.ResponseWriter, r *http.Request, status int, toRender interface{}, _ interface{}) error {
	if w == nil {
		return errors.New("render: http.ResponseWriter cannot be nil")
	}

	if r == nil {
		return errors.New("render: *http.Request cannot be nil")
	}

	var toRenderResponse interface{}

	switch toRenderType := toRender.(type) {
	case error:
		toRenderResponse = struct {
			Message string `json:"message"`
		}{
			toRenderType.Error(),
		}
	}

	js, err := json.Marshal(toRenderResponse)
	if err != nil {
		return errors.WithStack(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	return errors.WithStack(err)
}
