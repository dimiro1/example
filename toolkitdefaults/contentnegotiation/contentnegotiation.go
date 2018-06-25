package contentnegotiation

import (
	"net/http"
	"strings"

	"mime"

	"github.com/dimiro1/example/toolkitdefaults/contentnegotiation/internal/httputil"
)

// Negotiator ...
type Negotiator struct {
	// default: _format
	parameterName string
	// default application/json
	defaultType string

	offers []string
}

// Option ...
type Option func(*Negotiator)

// ParameterName option function to change the parameterName
//noinspection GoUnusedExportedFunction
func ParameterName(parameterName string) Option {
	return Option(func(n *Negotiator) { n.parameterName = parameterName })
}

// DefaultType option function to change the parameterName
func DefaultType(mime string) Option {
	return Option(func(n *Negotiator) { n.defaultType = mime })
}

func Offers(offers ...string) Option {
	return Option(func(n *Negotiator) { n.offers = offers })
}

// Negotiate basic implementation, first check the parameter on the querystring and after the Accept header
func (n *Negotiator) Negotiate(r *http.Request) string {
	ext := r.URL.Query().Get(n.parameterName)
	if len(ext) != 0 {
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		// Should we return */* ?
		return mime.TypeByExtension(ext)
	}

	bestOffer := httputil.NegotiateContentType(r, n.offers, n.defaultType)

	return bestOffer
}

func NewNegotiator(options ...Option) *Negotiator {
	n := &Negotiator{
		parameterName: "_format",
		defaultType:   "application/json",
		offers: []string{
			"application/json",
			"application/xml",
			"text/xml",
			"text/html",
			"application/xhtml+xml",
			"application/rss+xml",
			"application/atom+xml",
			"image/jpeg",
			"image/png",
			"image/gif",
			"text/markdown",
			"text/plain",
			"application/pdf",
		},
	}

	for _, option := range options {
		option(n)
	}

	return n
}
