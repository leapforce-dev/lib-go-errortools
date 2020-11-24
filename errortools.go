package errortools

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/getsentry/sentry-go"
)

// Println prints error if not nil
//
func Println(prefix string, err error) {
	if err != nil {
		fmt.Println(prefix, err)
	}
}

// Fatal prints error and exits if not nil
//
func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func captureError(err interface{}, toSentry bool) *Error {
	if !reflect.ValueOf(err).IsNil() {
		e := new(Error)

		if errError, ok := err.(*Error); ok {
			e = errError
		} else if errError, ok := err.(error); ok {
			e = ErrorMessage(errError)
		} else {
			//e = ErrorMessage(fmt.Sprintf("%v", err))
			e = ErrorMessage(nil)
		}

		if toSentry {
			if e.Response != nil {
				setTag("response_status_code", e.Response.StatusCode)
				setContext("response_status", e.Response.Status)
			} else {
				removeTag("response_status_code")
				removeContext("response_status")
			}

			if e.Request != nil {
				setContext("url", e.Request.URL.String())
				setContext("http_method", e.Request.Method)

				if e.Request.Body != nil {
					b := []byte{}
					_, _ = e.Request.Body.Read(b)
					setContext("http_body", b)
				} else {
					removeContext("http_body")
				}

			} else {
				removeContext("url")
				removeContext("http_method")
				removeContext("http_body")
			}

		}

		return e
	}

	return nil
}

// CaptureException sends error to Sentry, prints it and exits if not nil
//
func CaptureException(err interface{}, toSentry bool) {
	if !reflect.ValueOf(err).IsNil() {

		e := captureError(err, toSentry)
		if toSentry {
			sentry.CaptureException(errors.New(e.Message))
		}
		log.Fatal(err)
	}
}

// CaptureMessage sends message to Sentry, prints it and exits if not nil
//
func CaptureMessage(err *Error, toSentry bool) {

	if err != nil {
		e := captureError(err, toSentry)
		if toSentry {
			sentry.CaptureMessage(e.Message)
		}
		fmt.Println(e.Message)
	}
}

func setTag(key string, value interface{}) {
	sentry.CurrentHub().Scope().SetTag(key, fmt.Sprintf("%v", value))
}

func removeTag(key string) {
	sentry.CurrentHub().Scope().RemoveTag(key)
}

func setContext(key string, value interface{}) {
	sentry.CurrentHub().Scope().SetContext(key, fmt.Sprintf("%v", value))
}

func removeContext(key string) {
	sentry.CurrentHub().Scope().RemoveContext(key)
}
