package render

import (
	"encoding/json"
	"net/http"
)

type JSON struct{}

func (h JSON) Render(w http.ResponseWriter, r *http.Request, status int, toRender interface{}) error {
	return h.RenderCtx(w, r, status, toRender, nil)
}

func (JSON) RenderCtx(w http.ResponseWriter, r *http.Request, status int, toRender interface{}, context interface{}) error {
	switch toRender.(type) {
	case error:
		toRender = struct {
			Message string `json:"message"`
		}{
			toRender.(error).Error(),
		}
	}

	js, err := json.Marshal(toRender)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	return err
}
