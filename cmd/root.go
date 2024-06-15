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

type options struct {
	prompt     string
	promptPath string

	verbose bool
	silent  bool
	diff    bool

	rewrite bool
	outpath string
}

var opts options

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&opts.prompt, "prompt", "p", "", "Prompt text")
	rootCmd.Flags().StringVarP(&opts.promptPath, "prompt-path", "P", "", "Prompt file path")

	rootCmd.Flags().BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose mode")
	rootCmd.Flags().BoolVarP(&opts.silent, "silent", "s", false, "Silent mode")
	rootCmd.Flags().BoolVarP(&opts.diff, "diff", "d", false, "Show diff")

	rootCmd.Flags().BoolVarP(&opts.rewrite, "rewrite", "r", false, "Rewrite the input file with the result")
	rootCmd.Flags().StringVarP(&opts.outpath, "outpath", "o", "", "Output file path")
}

var rootCmd = &cobra.Command{
	Use:   "ai-text-shaper",
	Short: "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model",
	Long:  "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model.",
	// Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		verboseLog("ai-text-shaper started")
		verboseLog("opts: %+v", opts)
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

		outpath := opts.outpath
		if opts.rewrite {
			outpath = inputFilePath
		}

		verboseLog("start reading prompt")
		promptText, err := getPromptText(opts.prompt, opts.promptPath)
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

		if !opts.silent {
			if opts.diff {
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
