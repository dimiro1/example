package render

import (
	`encoding/xml`
	`net/http`
)

type XML struct{}

func (XML) Render(w http.ResponseWriter, r *http.Request, status int, i interface{}) error {
	switch i.(type) {
	case error:
		i = struct {
			XMLName xml.Name `xml:"error"`
			Message string   `xml:"message,attr"`
		}{
			Message: i.(error).Error(),
		}
	}

	x, err := xml.Marshal(i)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(status)
	_, err = w.Write(x)
	return err
}
