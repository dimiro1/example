package binder

import (
	"net/http"
	"github.com/dimiro1/example/internal/contenttype"
)

type ContentNegotiation struct {
	JSON       JSON
	XML        XML
	Parameters Gorilla
}

func (c ContentNegotiation) Bind(r *http.Request, dst interface{}) error {
	if r.ContentLength == 0 {
		if r.Method == http.MethodGet || r.Method == http.MethodDelete {
			return c.Parameters.Bind(r, dst)
		}
	}

	switch contenttype.Detect(r) {
	case "xml":
		return c.XML.Bind(r, dst)
	case "json":
		fallthrough
	default:
		return c.JSON.Bind(r, dst)
	}
}

func NewContentNegotiation() ContentNegotiation {
	return ContentNegotiation{}
}
