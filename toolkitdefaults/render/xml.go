package render

import (
	"encoding/xml"
	"net/http"

	"github.com/pkg/errors"
)

type XML struct{}

func (XML) Render(w http.ResponseWriter, r *http.Request, status int, toRender interface{}, _ interface{}) error {
	if w == nil {
		return errors.New("render: http.ResponseWriter cannot be nil")
	}

	if r == nil {
		return errors.New("render: *http.Request cannot be nil")
	}
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
		return errors.WithStack(err)
	}

	w.Header().Set("Content-Type", "text/xml")
	w.WriteHeader(status)
	_, err = w.Write(x)
	return errors.WithStack(err)
}
