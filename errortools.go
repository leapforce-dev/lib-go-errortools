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
			setExtra("error", e.message)
		} else {
			removeExtra("error")
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
func CaptureException(err interface{}, toSentry bool) {
	f, e := captureError(err, toSentry)
	if e != nil {
		if toSentry {
			sentry.CaptureException(errors.New(e.message))
		}
		log.Fatal(e.message)
	}

	if f != nil {
		f()
	}
}

// CaptureMessage sends message to Sentry, prints it and exits if not nil
//
func CaptureMessage(err interface{}, toSentry bool) {
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

func SetContext(key string, value string) {
	if context == nil {
		context = make(map[string]string)
	}
	context[key] = value
}

func RemoveContext(key string) {
	delete(context, key)
	//sentry.CurrentHub().Scope().removeExtra("context")
}
