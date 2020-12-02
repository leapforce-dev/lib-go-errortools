package errortools

import (
	"fmt"
	"net/http"
	"reflect"
)

// Error stores enriched error information
//
type Error struct {
	request  *http.Request
	response *http.Response
	message  string
	extras   *map[string]string
}

// ErrorMessage return message-only Error
//
func ErrorMessage(stringOrError interface{}) *Error {
	return &Error{message: message(stringOrError)}
}

func message(stringOrError interface{}) string {
	if errError, ok := stringOrError.(error); ok {
		return errError.Error()
	}
	if errString, ok := stringOrError.(string); ok {
		return errString
	}
	if errString, ok := stringOrError.(*string); ok {
		return *errString
	}
	return fmt.Sprintf("Invalid type %s (%s) passed to ErrorMessage.", reflect.TypeOf(stringOrError).Kind(), reflect.TypeOf(stringOrError))
}

func (err *Error) SetRequest(request *http.Request) {
	(*err).request = request
}

func (err *Error) SetResponse(response *http.Response) {
	(*err).response = response
}

func (err *Error) SetMessage(stringOrError interface{}) {
	(*err).message = message(stringOrError)
}

func (err *Error) SetExtra(key string, value string) {
	if (*err).extras == nil {
		m := make(map[string]string)
		(*err).extras = &m
	}

	(*((*err).extras))[key] = value
}

func (err *Error) Request() *http.Request {
	return err.request
}

func (err *Error) Response() *http.Response {
	return err.response
}

func (err *Error) Message() string {
	return err.message
}
