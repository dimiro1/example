package contenttype

import (
	"net/http"
)

const (
	JSON  = "json"
	XML   = "xml"
	PLAIN = "text"
	HTML  = "html"

	// more ...
)

// Detector ...
type Detector interface {
	Detect(r *http.Request, f func(theType string))
}
