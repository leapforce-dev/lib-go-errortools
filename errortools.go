package errortools

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/getsentry/sentry-go"
)

var context map[string]string

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

func captureError(err interface{}, toSentry bool) (func(), *Error) {
	if err == nil {
		return nil, nil
	}

	if reflect.TypeOf(err).Kind() == reflect.Ptr {
		if reflect.ValueOf(err).IsNil() {
			return nil, nil
		}
	}

	e := new(Error)

	if errError, ok := err.(*Error); ok {
		e = errError
	} else if errError, ok := err.(error); ok {
		e = ErrorMessage(errError)
	} else {
		e = ErrorMessage(fmt.Sprintf("%s: %v", reflect.TypeOf(err).String(), err))
	}

	removeFunc := func() {}

	if toSentry {
		if context != nil {
			c := []string{}

			for k, v := range context {
				c = append(c, fmt.Sprintf("%s: %s", k, v))
			}

			setExtra("context", strings.Join(c, "\n"))
		} else {
			removeExtra("context")
		}

		if e.message != "" {
			setExtra("message", e.message)
		} else {
			removeExtra("message")
		}

		if e.response != nil {
			setTag("response_status_code", e.response.StatusCode)
			setExtra("response_status", e.response.Status)
		} else {
			removeTag("response_status_code")
			removeExtra("response_status")
		}

		if e.request != nil {
			setExtra("url", e.request.URL.String())
			setExtra("http_method", e.request.Method)

			if e.request.Body != nil {
				readCloser, err := e.request.GetBody()
				if err != nil {
					fmt.Println(err)
				}
				b, err := ioutil.ReadAll(readCloser)
				if err == nil {
					setExtra("http_body", fmt.Sprintf("%s", b))
				} else {
					setExtra("http_body", fmt.Sprintf("Error reading body: %s", err.Error()))
				}
			} else {
				removeExtra("http_body")
			}

		} else {
			removeExtra("url")
			removeExtra("http_method")
			removeExtra("http_body")
		}

		if e.extras != nil {
			for key, value := range *(e.extras) {
				setExtra(key, value)
			}

			removeFunc = func() {
				for key, value := range *(e.extras) {
					setExtra(key, value)
				}
			}
		}
	}

	return removeFunc, e
}

// CaptureException sends error to Sentry, prints it and exits if not nil
//
func captureException(err interface{}, level sentry.Level, toSentry bool) {
	sentry.CurrentHub().Scope().SetLevel(level)
	defer sentry.CurrentHub().Scope().SetLevel("")

	f, e := captureError(err, toSentry)
	if e != nil {
		if toSentry {
			sentry.CaptureException(errors.New(e.message))
		}
		if level == sentry.LevelFatal {
			log.Fatal(e.message)
		} else {
			fmt.Println(e.message)
		}
	}

	if f != nil {
		f()
	}
}

// captureMessage sends message to Sentry, prints it and exits if not nil
//
func captureMessage(err interface{}, level sentry.Level, toSentry bool) {
	sentry.CurrentHub().Scope().SetLevel(level)
	defer sentry.CurrentHub().Scope().SetLevel("")

	f, e := captureError(err, toSentry)
	if e != nil {
		if toSentry {
			sentry.CaptureMessage(e.message)
		}
		fmt.Println(e.message)
	}

	if f != nil {
		f()
	}
}

// CaptureInfo sends info to Sentry, prints it and exits if not nil
//
func CaptureInfo(err interface{}, toSentry bool) {
	captureMessage(err, sentry.LevelInfo, toSentry)
}

// CaptureWarning sends warning to Sentry, prints it
//
func CaptureWarning(err interface{}, toSentry bool) {
	captureMessage(err, sentry.LevelWarning, toSentry)
}

// CaptureError sends error to Sentry, prints it
//
func CaptureError(err interface{}, toSentry bool) {
	captureException(err, sentry.LevelError, toSentry)
}

// CaptureFatal sends fatal to Sentry, prints it and exits if not nil
//
func CaptureFatal(err interface{}, toSentry bool) {
	captureException(err, sentry.LevelFatal, toSentry)
}

func setTag(key string, value interface{}) {
	sentry.CurrentHub().Scope().SetTag(key, fmt.Sprintf("%v", value))
}

func removeTag(key string) {
	sentry.CurrentHub().Scope().RemoveTag(key)
}

func setExtra(key string, value interface{}) {
	sentry.CurrentHub().Scope().SetExtra(key, fmt.Sprintf("%v", value))
}

func removeExtra(key string) {
	sentry.CurrentHub().Scope().RemoveExtra(key)
}

func SetContext(key string, value interface{}) {
	if context == nil {
		context = make(map[string]string)
	}
	context[key] = fmt.Sprintf("%v", value)
}

func RemoveContext(key string) {
	delete(context, key)
	//sentry.CurrentHub().Scope().removeExtra("context")
}
