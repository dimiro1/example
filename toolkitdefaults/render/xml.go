package render

import (
	`encoding/xml`
	`net/http`
)

type XML struct{}

func (XML) Render(w http.ResponseWriter, status int, toRender interface{}, context ...interface{}) error {
	switch toRender.(type) {
	case error:
		toRender = struct {
			XMLName xml.Name `xml:"error"`
			Message string   `xml:"message,attr"`
		}{
			Message: toRender.(error).Error(),
		}
	}

	x, err := xml.Marshal(toRender)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(status)
	_, err = w.Write(x)
	return err
}
