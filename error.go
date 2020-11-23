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

// ErrorMessage return message-only Error
//
func ErrorMessage(stringOrError interface{}) *Error {
	errError, ok := stringOrError.(error)
	if ok {
		return &Error{Message: errError.Error()}
	}
	errString, ok := stringOrError.(string)
	if ok {
		return &Error{Message: errString}
	}
	return &Error{Message: "Invalid type passed to ErrorMessage."}
}
