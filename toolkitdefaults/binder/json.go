package binder

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type JSON struct{}

func (JSON) Bind(r *http.Request, dst interface{}) error {
	if r == nil {
		return errors.New("render: *http.Request cannot be nil")
	}
	return errors.WithStack(json.NewDecoder(r.Body).Decode(dst))
}
