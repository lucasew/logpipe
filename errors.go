package logpipe

import "log"

// ReportError is the centralized error reporting function.
func ReportError(err error) {
	if err != nil {
		log.Printf("error: %v", err)
	}
}
