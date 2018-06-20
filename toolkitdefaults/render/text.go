package render

import (
	`fmt`
	`net/http`
)

type Text struct{}

func (Text) Render(w http.ResponseWriter, status int, i interface{}, extra ...interface{}) error {
	w.Header().Set("Content-Type", "text/plain")
	var data []byte

	// Specific types
	switch i.(type) {
	case string:
		data = []byte(i.(string))
	case error:
		data = []byte(i.(error).Error())
	}

	// Stringer
	if s, ok := i.(fmt.Stringer); ok {
		data = []byte(s.String())
	}

	w.WriteHeader(status)
	_, err := w.Write(data)
	return err
}
