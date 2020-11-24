package errortools

import (
	"fmt"
	"net/http"
	"reflect"
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
	if errError, ok := stringOrError.(error); ok {
		return errError.Error()
	}
	if errString, ok := stringOrError.(string); ok {
		return errString
	}
	return fmt.Sprintf("Invalid type %s passed to ErrorMessage.", reflect.TypeOf(stringOrError).Kind())
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
