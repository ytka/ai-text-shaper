package cmd

import "log"

func verboseLog(msg string, args ...interface{}) {
	if fl.verbose {
		log.Printf(msg, args...)
	}
}
