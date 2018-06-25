package recipes

import (
	"encoding/xml"
)

type errorResponse struct {
	XMLName xml.Name `json:"-" xml:"error"`
	Message string   `json:"message" xml:"message,attr"`
}
