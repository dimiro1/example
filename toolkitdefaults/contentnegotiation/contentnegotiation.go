package contentnegotiation

import (
	"net/http"
	"strings"

	"github.com/dimiro1/example/toolkit/contentnegotiation/mediatype"
)

// Negotiator ...
type Negotiator struct {
	// default: mediaType
	parameterName string
	// default application/json
	mediaType mediatype.MediaType
}

// Option ...
type Option func(*Negotiator)

// ParameterName option function to change the parameterName
func ParameterName(name string) Option {
	return Option(func(n *Negotiator) { n.parameterName = name })
}

// MediaType option function to change the parameterName
func MediaType(mediaType mediatype.MediaType) Option {
	return Option(func(n *Negotiator) { n.mediaType = mediaType })
}

// Negotiate basic implementation, first check the parameter on the querystring and after the Accept header
// TODO: Deal with priorities in the Accept header
func (n *Negotiator) Negotiate(r *http.Request) mediatype.MediaType {
	// The format parameter
	ext := map[string]mediatype.MediaType{
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
	types := map[string]mediatype.MediaType{
		mediatype.ApplicationAtomXML.String():        mediatype.ApplicationAtomXML,
		mediatype.ApplicationFormURLEncoded.String(): mediatype.ApplicationFormURLEncoded,
		mediatype.ApplicationPDF.String():            mediatype.ApplicationPDF,
		mediatype.ApplicationOctetStream.String():    mediatype.ApplicationOctetStream,
		mediatype.ApplicationJSON.String():           mediatype.ApplicationJSON,
		mediatype.ApplicationRSSXML.String():         mediatype.ApplicationRSSXML,
		mediatype.ApplicationXHTMLXML.String():       mediatype.ApplicationXHTMLXML,
		mediatype.ApplicationXML.String():            mediatype.ApplicationXML,
		mediatype.ImageGif.String():                  mediatype.ImageGif,
		mediatype.ImageJpeg.String():                 mediatype.ImageJpeg,
		mediatype.ImagePng.String():                  mediatype.ImagePng,
		mediatype.TextHTML.String():                  mediatype.TextHTML,
		mediatype.TextMarkdown.String():              mediatype.TextMarkdown,
		mediatype.TextPlain.String():                 mediatype.TextPlain,
		mediatype.TextXML.String():                   mediatype.ApplicationXML, // Special case
	}

	// return the default
	if strings.Contains(acceptHeader, mediatype.All.String()) {
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
