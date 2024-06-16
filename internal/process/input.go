package process

import (
	"fmt"
	"io"
	"os"
)

func GetInputText(inputFilePath string) (string, error) {
	if inputFilePath == "-" {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("error reading input from stdin: %w", err)
		}
		return string(input), nil
	}

	input, err := os.ReadFile(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("error reading input file: %w", err)
	}
	return string(input), nil
}
