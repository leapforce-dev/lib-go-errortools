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
	return &Error{Message: message(stringOrError)}
}

func message(stringOrError interface{}) string {
	if stringOrError == nil {
		return ""
	}
	errError, ok := stringOrError.(error)
	if ok {
		return errError.Error()
	}
	errString, ok := stringOrError.(string)
	if ok {
		return errString
	}
	return "Invalid type passed to ErrorMessage."
}

func (err *Error) SetRequest(request *http.Request) {
	(*err).Request = request
}

func (err *Error) SetResponse(response *http.Response) {
	(*err).Response = response
}

func (err *Error) SetMessage(stringOrError interface{}) {
	(*err).Message = message(stringOrError)
}
