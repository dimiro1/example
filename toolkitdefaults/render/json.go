package render

import (
	`encoding/json`
	`net/http`
)

type JSON struct{}

func (JSON) Render(w http.ResponseWriter, status int, i interface{}, extra ...interface{}) error {
	switch i.(type) {
	case error:
		i = struct {
			Message string `json:"message"`
		}{
			i.(error).Error(),
		}
	}

	js, err := json.Marshal(i)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(js)
	return err
}
