package render

import (
	`net/http`

	"github.com/dimiro1/example/internal/contenttype"
)

type ContentNegotiation struct {
	JSONRenderer JSON
	XMLRenderer  XML
	TextRenderer Text
}

func (c ContentNegotiation) Render(w http.ResponseWriter, r *http.Request, status int, i interface{}) error {
	switch contenttype.Detect(r) {
	case "xml":
		return c.XMLRenderer.Render(w, r, status, i)
	case "json":
		return c.JSONRenderer.Render(w, r, status, i)
	default:
		return c.TextRenderer.Render(w, r, status, i)
	}
}

func NewContentNegotiation() ContentNegotiation {
	return ContentNegotiation{}
}
