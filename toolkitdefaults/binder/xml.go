package binder

import (
	"encoding/xml"
	"net/http"

	"github.com/pkg/errors"
)

type XML struct{}

func (XML) Bind(r *http.Request, dst interface{}) error {
	if r == nil {
		return errors.New("render: *http.Request cannot be nil")
	}
	return errors.WithStack(xml.NewDecoder(r.Body).Decode(dst))
}
