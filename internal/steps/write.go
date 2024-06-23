package steps

import (
	"fmt"
	"os"
)

// WriteResult writes the outputText to the given outpath. If an error occurs, it wraps it with additional context.
func WriteResult(outputText, outpath string) error {
	if err := os.WriteFile(outpath, []byte(outputText), 0o600); err != nil {
		return fmt.Errorf("error writing to file: %w", err) // fixed: Error messages should start with lowercase letter and should not include punctuation. Refer to the section "Error Strings" in the Google Go style guide.
	}
	return nil
}
