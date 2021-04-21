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
var modifyMessageFunction *func(message string) string
var errorCount int

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

func captureError(err interface{}) (func(), *Error) {
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
		e.originalError = errError
	} else if errError, ok := err.(string); ok {
		e = ErrorMessage(errError)
	} else if errError, ok := err.(*string); ok {
		e = ErrorMessage(errError)
	} else {
		e = ErrorMessage(fmt.Sprintf("%s: %v", reflect.TypeOf(err).String(), err))
	}

	if modifyMessageFunction != nil {
		e.SetMessage((*modifyMessageFunction)(e.Message()))
	}

	removeFunc := func() {}

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

	return removeFunc, e
}

// CaptureException sends error to Sentry, prints it and exits if not nil
//
func captureException(err interface{}, level sentry.Level) {

	f, e := captureError(err)

	if e != nil {
		sentry.CurrentHub().Scope().SetLevel(level)
		defer sentry.CurrentHub().Scope().SetLevel("")

		if e.originalError != nil {
			sentry.CaptureException(e.originalError)
		} else {
			sentry.CaptureException(errors.New(e.message))
		}

		if level == sentry.LevelFatal {
			log.Fatal(e.message)
		} else {
			errorCount++
			fmt.Println(e.message)
		}
	}

	if f != nil {
		f()
	}
}

// captureMessage sends message to Sentry, prints it and exits if not nil
//
func captureMessage(err interface{}, level sentry.Level) {

	f, e := captureError(err)
	if e != nil {
		sentry.CurrentHub().Scope().SetLevel(level)
		defer sentry.CurrentHub().Scope().SetLevel("")
		sentry.CaptureMessage(e.message)

		fmt.Println(e.message)
	}

	if f != nil {
		f()
	}
}

// CaptureInfo sends info to Sentry, prints it and exits if not nil
//
func CaptureInfo(err interface{}) {
	captureMessage(err, sentry.LevelInfo)
}

// CaptureInfof returns formatted message-only Error
//
func CaptureInfof(format string, a ...interface{}) {
	CaptureInfo(fmt.Sprintf(format, a...))
}

// CaptureWarning sends warning to Sentry, prints it
//
func CaptureWarning(err interface{}) {
	captureMessage(err, sentry.LevelWarning)
}

// CaptureWarningf returns formatted message-only Error
//
func CaptureWarningf(format string, a ...interface{}) {
	CaptureWarning(fmt.Sprintf(format, a...))
}

// CaptureError sends error to Sentry, prints it
//
func CaptureError(err interface{}) {
	captureException(err, sentry.LevelError)
}

// CaptureErrorf returns formatted message-only Error
//
func CaptureErrorf(format string, a ...interface{}) {
	CaptureError(fmt.Sprintf(format, a...))
}

// CaptureFatal sends fatal to Sentry, prints it and exits if not nil
//
func CaptureFatal(err interface{}) {
	captureException(err, sentry.LevelFatal)
}

// CaptureFatalf returns formatted message-only Error
//
func CaptureFatalf(format string, a ...interface{}) {
	CaptureFatal(fmt.Sprintf(format, a...))
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

func SetModifyMessageFunction(function *func(message string) string) {
	modifyMessageFunction = function
}

func RemoveModifyMessageFunction() {
	modifyMessageFunction = nil
}

func Count() int {
	return errorCount
}

func ResetCount() {
	errorCount = 0
}
