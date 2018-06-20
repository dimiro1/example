package binder

import (
	"encoding/xml"
	"net/http"
)

type XML struct{}

func (XML) Bind(r *http.Request, dst interface{}) error {
	return xml.NewDecoder(r.Body).Decode(dst)
}
