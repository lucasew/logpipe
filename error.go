package logpipe

import (
	"log"
)

// ReportError centralizes error reporting for the application.
// This function ensures all unexpected errors are logged consistently
// and avoids scattered panics or empty catch blocks.
func ReportError(err error) {
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		// Here, we could hook into Sentry or other error reporting backends.
	}
}
