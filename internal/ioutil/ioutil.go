package ioutil

import (
	"errors"
	"fmt"
	"os"
)

// ErrFileIsNil is an error for a nil file.
var ErrFileIsNil = errors.New("file is nil")

// IsAvailablePipe checks whether the given file is a named pipe and has content.
func IsAvailablePipe(file *os.File) (bool, error) {
	if file == nil {
		return false, ErrFileIsNil
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return false, fmt.Errorf("failed to get file info: %w", err)
	}
	return fileInfo.Mode()&os.ModeNamedPipe != 0 && fileInfo.Size() > 0, nil
}
