package errortools

import (
	"time"

	"github.com/getsentry/sentry-go"
)

func InitSentry(dsn string) {
	// We need to use the sync transport (which is not the default),
	// otherwise if you use "log.Fatal()" the program will exit before the
	// error is sent to Sentry.
	sentrySyncTransport := sentry.NewHTTPSyncTransport()
	sentrySyncTransport.Timeout = time.Second * 3
	err := sentry.Init(sentry.ClientOptions{
		AttachStacktrace: true,
		Dsn:              dsn,
		Transport:        sentrySyncTransport,
	})
	if err != nil {
		panic(err)
	}
}
