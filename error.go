package errortools

import (
	"fmt"
	"net/http"
	"reflect"
)

// Error stores enriched error information
//
type Error struct {
	originalError error
	request       *http.Request
	body          []byte
	response      *http.Response
	message       string
	fingerprint   *[]string
	extras        *map[string]string
}

// ErrorMessage returns message-only Error
//
func ErrorMessage(stringOrError interface{}) *Error {
	return &Error{message: message(stringOrError)}
}

// ErrorMessagef returns formatted message-only Error
//
func ErrorMessagef(format string, a ...interface{}) *Error {
	return &Error{message: fmt.Sprintf(format, a...)}
}

func message(stringOrError interface{}) string {
	if errError, ok := stringOrError.(*Error); ok {
		return errError.Message()
	}
	if errError, ok := stringOrError.(error); ok {
		return errError.Error()
	}
	if errString, ok := stringOrError.(string); ok {
		return errString
	}
	if errString, ok := stringOrError.(*string); ok {
		return *errString
	}
	return fmt.Sprintf("Invalid type %s (%s) passed to message.", reflect.TypeOf(stringOrError).Kind(), reflect.TypeOf(stringOrError))
}

func (err *Error) SetRequest(request *http.Request) {
	(*err).request = request
}

func (err *Error) SetBody(b []byte) {
	(*err).body = b
}

func (err *Error) SetResponse(response *http.Response) {
	(*err).response = response
}

func (err *Error) SetMessage(stringOrError interface{}) {
	(*err).message = message(stringOrError)
}

func (err *Error) SetMessagef(format string, a ...interface{}) {
	(*err).message = message(fmt.Sprintf(format, a...))
}

func (err *Error) SetFingerprint(fingerprint *[]string) {
	(*err).fingerprint = fingerprint
}

func (err *Error) SetExtra(key string, value string) {
	if (*err).extras == nil {
		m := make(map[string]string)
		(*err).extras = &m
	}

	(*((*err).extras))[key] = value
}

func (err *Error) SetType(value string) {
	err.SetExtra(KeyExceptionType, value)
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
