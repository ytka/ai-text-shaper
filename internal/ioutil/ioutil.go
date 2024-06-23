package ioutil

import (
	"fmt"
	"os"
)

func IsAvailablePipe(file *os.File) (bool, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return false, fmt.Errorf("failed to get file info: %w", err)
	}
	return fileInfo.Mode()&os.ModeNamedPipe != 0 && fileInfo.Size() > 0, nil
}
