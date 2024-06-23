package steps

import (
	"fmt"
	"io"
	"os"
)

// getInputFromStdin reads all input from stdin if not connected to a tty.
func getInputFromStdin() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("error getting stdin stat: %w", err)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", nil
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("error reading input from stdin: %w", err)
	}
	return string(input), nil
}

// GetInputText reads text from a file specified by inputFilePath or from stdin if the path is "-".
func GetInputText(inputFilePath string) (string, error) { // fixed: Function name should follow MixedCaps style - https://google.github.io/styleguide/go/guide.html#mixed-caps
	if inputFilePath == "-" {
		return getInputFromStdin()
	}

	input, err := os.ReadFile(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("error reading input file: %w", err)
	}
	return string(input), nil
}
