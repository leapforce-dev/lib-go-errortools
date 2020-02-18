package errortools

import (
	"fmt"
	"log"

	sentry "github.com/getsentry/sentry-go"
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

// FatalSentry sends error to Sentry, prints it and exits if not nil
//
func FatalSentry(err error, isLive bool) {
	if err != nil {
		if isLive {
			sentry.CaptureException(err)
		}
		log.Fatal(err)
	}
}
