package contentnegotiation

import (
	"net/http"
	"strings"

	"github.com/dimiro1/example/toolkit/mediatype"
)

// Negotiator ...
type Negotiator struct {
	// default: mediaType
	parameterName string
	// default application/json
	mediaType string
}

// Option ...
type Option func(*Negotiator)

// ParameterName option function to change the parameterName
func ParameterName(name string) Option {
	return Option(func(n *Negotiator) { n.parameterName = name })
}

// MediaType option function to change the parameterName
func MediaType(mediaType string) Option {
	return Option(func(n *Negotiator) { n.mediaType = mediaType })
}

// Negotiate basic implementation, first check the parameter on the querystring and after the Accept header
// TODO: Deal with priorities in the Accept header
func (n *Negotiator) Negotiate(r *http.Request) string {
	// The format parameter
	ext := map[string]string{
		"atom":  mediatype.ApplicationAtomXML,
		"pdf":   mediatype.ApplicationPDF,
		"json":  mediatype.ApplicationJSON,
		"rss":   mediatype.ApplicationRSSXML,
		"xhtml": mediatype.ApplicationXHTMLXML,
		"xml":   mediatype.ApplicationXML,
		"gif":   mediatype.ImageGif,
		"jpeg":  mediatype.ImageJpeg,
		"jpg":   mediatype.ImageJpeg,
		"png":   mediatype.ImagePng,
		"txt":   mediatype.TextHTML,
		"md":    mediatype.TextMarkdown,
		"html":  mediatype.TextHTML,
	}

	parameter := r.URL.Query().Get(n.parameterName)
	for value, tType := range ext {
		if strings.HasPrefix(parameter, value) {
			return tType
		}
	}

	// Accept header
	acceptHeader := r.Header.Get("Accept")
	types := map[string]string{
		mediatype.ApplicationAtomXML:        mediatype.ApplicationAtomXML,
		mediatype.ApplicationFormURLEncoded: mediatype.ApplicationFormURLEncoded,
		mediatype.ApplicationPDF:            mediatype.ApplicationPDF,
		mediatype.ApplicationOctetStream:    mediatype.ApplicationOctetStream,
		mediatype.ApplicationJSON:           mediatype.ApplicationJSON,
		mediatype.ApplicationRSSXML:         mediatype.ApplicationRSSXML,
		mediatype.ApplicationXHTMLXML:       mediatype.ApplicationXHTMLXML,
		mediatype.ApplicationXML:            mediatype.ApplicationXML,
		mediatype.ImageGif:                  mediatype.ImageGif,
		mediatype.ImageJpeg:                 mediatype.ImageJpeg,
		mediatype.ImagePng:                  mediatype.ImagePng,
		mediatype.TextHTML:                  mediatype.TextHTML,
		mediatype.TextMarkdown:              mediatype.TextMarkdown,
		mediatype.TextPlain:                 mediatype.TextPlain,
		mediatype.TextXML:                   mediatype.ApplicationXML, // Special case
	}

	// return the default
	if strings.Contains(acceptHeader, mediatype.All) {
		return n.mediaType
	}

	for value, tType := range types {
		if strings.HasPrefix(acceptHeader, value) {
			return tType
		}
	}

	return n.mediaType
}

func NewNegotiator(options ...Option) *Negotiator {
	n := &Negotiator{"mediaType", mediatype.ApplicationJSON}

	for _, option := range options {
		option(n)
	}

	return n
}
