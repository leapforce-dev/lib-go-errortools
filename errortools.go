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

func captureError(err *Error, toSentry bool) {
	if err != nil {
		if toSentry {
			if err.Response != nil {
				setTag("ResponseStatusCode", err.Response.StatusCode)
			} else {
				removeTag("ResponseStatusCode")
			}

			if err.Request != nil {
				b := []byte{}
				_, _ = err.Request.Body.Read(b)

				setContext("Body", b)
			} else {
				removeContext("Body")
			}

		}
	}
}

// CaptureException sends error to Sentry, prints it and exits if not nil
//
func CaptureException(err *Error, toSentry bool) {

	if err != nil {
		captureError(err, toSentry)
		if toSentry {
			sentry.CaptureException(errors.New(err.Message))
		}
		log.Fatal(err)
		fmt.Println(err.Message)
	}
}

// CaptureMessage sends message to Sentry, prints it and exits if not nil
//
func CaptureMessage(err *Error, toSentry bool) {

	if err != nil {
		captureError(err, toSentry)
		if toSentry {
			sentry.CaptureMessage(err.Message)
		}
		fmt.Println(err.Message)
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
