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

var fl flags

func init() {
	fl.initCommandFlags(rootCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ai-text-shaper",
	Short: "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model",
	Long:  "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model.",
	// Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		verboseLog("ai-text-shaper started")
		verboseLog("flags: %+v", fl)
		verboseLog("args: %v", args)

		apikey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		inputFilePath := "-"
		if len(args) >= 1 {
			// FIXME: larger case
			inputFilePath = args[0]
		}

		outpath := fl.outpath
		if fl.rewrite {
			outpath = inputFilePath
		}

		verboseLog("start reading prompt")
		promptText, err := getPromptText(fl.prompt, fl.promptPath)
		if err != nil {
			return err
		}

		verboseLog("start reading input: %s", inputFilePath)
		inputText, err := getInputText(inputFilePath)
		if err != nil {
			return err
		}

		verboseLog("start shaping text")
		resultText, err := textshaper.ShapeText(apikey, promptText, inputText)
		verboseLog("end shaping text")
		if err != nil {
			return err
		}

		if !fl.silent {
			if fl.diff {
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
