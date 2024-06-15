package main

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"os"
	"path/filepath"
	"strings"
)

func getAPIKey() (string, error) {
	apiKeyFilePath := os.Getenv("HOME") + "/.openai-apikey"
	apiKey, err := os.ReadFile(apiKeyFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}
	return string(apiKey), nil
}

func shapeText(apiKey, prompt, input string) (string, error) {
	mergedPrmpt := fmt.Sprintf("%s\n\n%s", prompt, input)
	result, err := callOpenAI(apiKey, "gpt-4o", mergedPrmpt)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	result = strings.TrimSuffix(result, "\n")
	result = strings.TrimSpace(result)
	lines := strings.Split(result, "\n")
	if len(lines) > 0 && strings.HasPrefix(result, "```") {
		lines = lines[1:]
	}
	if len(lines) > 0 && strings.HasSuffix(result, "```") {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n"), nil
}

func getFilePathWithoutExt(filePath string) string {
	return filePath[:len(filePath)-len(filepath.Ext(filePath))]
}

func writeToFile(filePath, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

// TODO: multiple input files
func main() {
	apiKey, err := getAPIKey()
	apiKey = strings.TrimSuffix(apiKey, "\n")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	args := os.Args
	if len(args) != 4 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s <command> <prompt file> <input file>\n", args[0])
		return
	}
	command := args[1] // "rewrite" | "diff" | "print" | "gen"
	promptFilePath := args[2]
	inputFilePath := args[3]
	promptName := filepath.Base(getFilePathWithoutExt(promptFilePath))

	prompt, err := os.ReadFile(promptFilePath)
	if err != nil {
		fmt.Println("Error reading prompt file:", err)
		return
	}
	input, err := os.ReadFile(inputFilePath)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}
	inputText := string(input)

	resultText, err := shapeText(apiKey, string(prompt), inputText)
	if err != nil {
		fmt.Println("Error shaping text:", err)
		return
	}

	switch command {
	case "print":
		fmt.Println(resultText)
	case "diff":
		dmp := diffmatchpatch.New()
		a, b, c := dmp.DiffLinesToChars(inputText, resultText)
		diffs := dmp.DiffMain(a, b, false)
		diffs = dmp.DiffCharsToLines(diffs, c)
		fmt.Println(dmp.DiffPrettyText(diffs))
	case "gen":
		outputFilePath := fmt.Sprintf("%s-%s%s", getFilePathWithoutExt(inputFilePath), promptName, filepath.Ext(inputFilePath))
		err = writeToFile(outputFilePath, resultText)
		if err != nil {
			fmt.Println("gen error:", err)
			return
		}
	case "rewrite":
		if err := writeToFile(inputFilePath, resultText); err != nil {
			fmt.Println("Error writing to file:", err)
		}
	default:
		fmt.Println("Unsupported command:", command)
	}
}
