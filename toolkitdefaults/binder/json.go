package binder

import (
	"encoding/json"
	"net/http"
)

type JSON struct{}

func (JSON) Bind(r *http.Request, dst interface{}) error {
	return json.NewDecoder(r.Body).Decode(dst)
}
