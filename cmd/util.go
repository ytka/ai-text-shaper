package cmd

import "log"

func verboseLog(msg string, args ...interface{}) {
	if opts.verbose {
		log.Printf(msg, args...)
	}
}
