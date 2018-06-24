package mediatype

type MediaType string

func (m MediaType) String() string {
	return string(m)
}

const (
	ApplicationAtomXML        MediaType = "application/atom+xml"
	ApplicationFormURLEncoded MediaType = "application/x-www-form-urlencoded"
	ApplicationPDF            MediaType = "application/pdf"
	ApplicationOctetStream    MediaType = "application/octet-stream"
	ApplicationJSON           MediaType = "application/json"
	ApplicationRSSXML         MediaType = "application/rss+xml"
	ApplicationXHTMLXML       MediaType = "application/xhtml+xml"
	ApplicationXML            MediaType = "application/xml"

	ImageGif  MediaType = "image/gif"
	ImageJpeg MediaType = "image/jpeg"
	ImagePng  MediaType = "image/png"

	TextHTML     MediaType = "text/html"
	TextMarkdown MediaType = "text/markdown"
	TextPlain    MediaType = "text/plain"
	TextXML      MediaType = "text/xml"

	All MediaType = "*/*"
)
