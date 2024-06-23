package steps

import (
	"fmt"
	"io"
	"os"
)

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

func GetInputText(inputFilePath string) (string, error) {
	if inputFilePath == "-" {
		return getInputFromStdin()
	}

	input, err := os.ReadFile(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("error reading input file: %w", err)
	}
	return string(input), nil
}
