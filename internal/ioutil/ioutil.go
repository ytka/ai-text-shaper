package ioutil

import (
	"errors"
	"fmt"
	"os"
)

// ErrFileIsNil is an error for a nil file.
var ErrFileIsNil = errors.New("file is nil")

// IsStdinPipe checks if stdin is a pipe.
func IsStdinPipe() (bool, error) {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false, fmt.Errorf("failed to get file info: %w", err)
	}
	return fileInfo.Mode()&os.ModeNamedPipe != 0, nil
}

// IsStdoutPipeOrRedirect checks if stdout is a pipe or a redirect.
func IsStdoutPipeOrRedirect() (bool, error) {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false, fmt.Errorf("failed to get file info: %w", err)
	}
	return fileInfo.Mode()&os.ModeCharDevice == 0, nil
}
