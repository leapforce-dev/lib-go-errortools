package errortools

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"

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

// CaptureException sends error to Sentry, prints it and exits if not nil
//
func CaptureException(errorOrString interface{}, responseStatusCode int, data interface{}, toSentry bool) {
	if errorOrString != nil {
		var err error

		errError, ok := errorOrString.(*error)
		if ok {
			err = *errError
		} else {
			errString, ok := errorOrString.(string)
			if ok {
				err = errors.New(errString)
			} else {
				err = fmt.Errorf("Error of type %s: %s", reflect.ValueOf(errorOrString).Kind(), fmt.Sprintf("%v", errorOrString))
			}
		}

		if toSentry {
			setTagAndContext(responseStatusCode, data)
			sentry.CaptureException(err)
		}
		log.Fatal(err)
	}
}

// CaptureMessage sends message to Sentry, prints it and exits if not nil
//
func CaptureMessage(message string, responseStatusCode int, data interface{}, toSentry bool) {
	if toSentry {
		setTagAndContext(responseStatusCode, data)
		sentry.CaptureMessage(message)
	}
	fmt.Println(message)
}

func setTagAndContext(responseStatusCode int, data interface{}) {
	if responseStatusCode > 0 {
		sentry.CurrentHub().Scope().SetTag("ResponseStatusCode", strconv.Itoa(responseStatusCode))
	} else {
		sentry.CurrentHub().Scope().RemoveTag("ResponseStatusCode")
	}
	if data != nil {
		sentry.CurrentHub().Scope().SetContext("Data", data)
	} else {
		sentry.CurrentHub().Scope().RemoveContext("Data")
	}
}
