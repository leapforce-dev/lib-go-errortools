package errortools

import (
	"errors"
	"fmt"
	"log"

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
	if err != nil {
		e := new(Error)

		if errError, ok := err.(*Error); ok {
			e = errError
		} else if errError, ok := err.(error); ok {
			e = ErrorMessage(errError)
		} else {
			e = ErrorMessage(fmt.Sprintf("%v", err))
		}

		if toSentry {
			if e.Response != nil {
				setTag("ResponseStatusCode", e.Response.StatusCode)
				setContext("ResponseStatus", e.Response.Status)
			} else {
				removeTag("ResponseStatusCode")
				removeContext("ResponseStatus")
			}

			if e.Request != nil {
				b := []byte{}
				_, _ = e.Request.Body.Read(b)

				setContext("URL", e.Request.URL.String())
				setContext("Method", e.Request.Method)
				setContext("Body", b)
			} else {
				removeContext("URL")
				removeContext("Method")
				removeContext("Body")
			}

		}

		return e
	}

	return nil
}

// CaptureException sends error to Sentry, prints it and exits if not nil
//
func CaptureException(err interface{}, toSentry bool) {

	if err != nil {
		e := captureError(err, toSentry)
		if toSentry {
			sentry.CaptureException(errors.New(e.Message))
		}
		log.Fatal(err)
		fmt.Println(e.Message)
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
