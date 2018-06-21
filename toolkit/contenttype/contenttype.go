package contenttype

import (
	"net/http"
)

const (
	JSON  = "json"
	XML   = "xml"
	PLAIN = "text"
	HTML  = "html"
	ANY   = "any"

	// more ...
)

// Detector ...
type Detector interface {
	Detect(r *http.Request) string
}
