package render

import (
	`encoding/xml`
	`net/http`
)

type XML struct{}

func (h XML) Render(w http.ResponseWriter, r *http.Request, status int, toRender interface{}) error {
	return h.RenderCtx(w, r, status, toRender, nil)
}

func (XML) RenderCtx(w http.ResponseWriter, r *http.Request, status int, toRender interface{}, context interface{}) error {
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
