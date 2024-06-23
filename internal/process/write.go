package process

import (
	"fmt"
	"os"
)

func WriteResult(outputText, outpath string) error {
	if err := os.WriteFile(outpath, []byte(outputText), 0o644); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}
