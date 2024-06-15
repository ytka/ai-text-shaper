package cmd

import (
	"ai-text-shaper/internal/textshaper"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var (
	prompt     string
	promptPath string

	silent bool
	diff   bool

	rewrite bool
	outpath string
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&prompt, "prompt", "p", "", "Prompt text")
	rootCmd.Flags().StringVarP(&promptPath, "prompt-path", "P", "", "Prompt file path")

	rootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "Silent mode")
	rootCmd.Flags().BoolVarP(&diff, "diff", "d", false, "Show diff")

	rootCmd.Flags().BoolVarP(&rewrite, "rewrite", "r", false, "Rewrite the input file with the result")
	rootCmd.Flags().StringVarP(&outpath, "outpath", "o", "", "Output file path")
}

var rootCmd = &cobra.Command{
	Use:   "ai-text-shaper",
	Short: "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model",
	Long:  "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model.",
	// Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apikey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		inputFilePath := "-"
		if len(args) > 1 {
			// FIXME: larger case
			inputFilePath = args[0]
		}

		if rewrite {
			outpath = inputFilePath
		}

		promptText, err := getPromptText(prompt, promptPath)
		if err != nil {
			return err
		}

		inputText, err := getInputText(inputFilePath)
		if err != nil {
			return err
		}

		resultText, err := textshaper.ShapeText(apikey, promptText, inputText)
		if err != nil {
			return err
		}

		if !silent {
			if diff {
				dmp := diffmatchpatch.New()
				a, b, c := dmp.DiffLinesToChars(inputText, resultText)
				diffs := dmp.DiffMain(a, b, false)
				diffs = dmp.DiffCharsToLines(diffs, c)
				fmt.Println(dmp.DiffPrettyText(diffs))
			} else {
				fmt.Println(resultText)
			}
		}
		if outpath != "" {
			fmt.Println("Writing to file...", outpath)
			if err := os.WriteFile(outpath, []byte(resultText), 0644); err != nil {
				return fmt.Errorf("error writing to file: %w", err)
			}
		}

		return nil
	},
}

func getAPIKey() (string, error) {
	apiKeyFilePath := os.Getenv("HOME") + "/.openai-apikey"
	bytes, err := os.ReadFile(apiKeyFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}
	return strings.TrimSuffix(string(bytes), "\n"), nil
}

func getPromptText(prompt, promptPath string) (string, error) {
	if prompt == "" && promptPath == "" {
		return "", fmt.Errorf("prompt is required")
	}
	if prompt == "" && promptPath != "" {
		text, err := os.ReadFile(promptPath)
		if err != nil {
			return "", fmt.Errorf("error reading prompt file: %w", err)
		}
		return string(text), nil
	}
	return prompt, nil
}

func getInputText(inputFilePath string) (string, error) {
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
