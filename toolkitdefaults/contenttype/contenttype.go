package contenttype

import (
	"net/http"
	internalContentType "github.com/dimiro1/example/internal/contenttype"
)

// Detector ...
type Detector struct{}

func (Detector) Detect(r *http.Request, f func(theType string)) {
	f(internalContentType.Detect(r))
}
