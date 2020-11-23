package errortools

import (
	"net/http"
)

// Error stores enriched error information
//
type Error struct {
	Request  *http.Request
	Response *http.Response
	Message  string
}
