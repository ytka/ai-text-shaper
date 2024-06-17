package cmd

import (
	"ai-text-shaper/internal/openai"
	"ai-text-shaper/internal/process"
	"ai-text-shaper/internal/runner"
	"ai-text-shaper/internal/tui"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var c runner.Config

func init() {
	rootCmd.Version = "testX"
	// prompt options
	rootCmd.Flags().StringVarP(&c.Prompt, "prompt", "p", "", "Prompt text")
	rootCmd.Flags().StringVarP(&c.PromptPath, "prompt-path", "P", "", "Prompt file path")

	// stdout messages options
	rootCmd.Flags().BoolVarP(&c.Verbose, "verbose", "v", false, "Verbose output")
	rootCmd.Flags().BoolVarP(&c.Silent, "silent", "s", false, "Suppress output")
	rootCmd.Flags().BoolVarP(&c.Diff, "diff", "d", false, "Show diff of the input and output text")

	// write file options
	rootCmd.Flags().BoolVarP(&c.Rewrite, "rewrite", "r", false, "Rewrite the input file with the result")
	rootCmd.Flags().StringVarP(&c.Outpath, "outpath", "o", "", "Output file path")
	rootCmd.Flags().BoolVarP(&c.UseFirstCodeBlock, "use-first-code-block", "f", false, "Use the first code block in the output text")
	rootCmd.Flags().BoolVarP(&c.Confirm, "confirm", "c", false, "Confirm before writing to file")

	// model options
	rootCmd.Flags().StringVarP(&c.Model, "model", "m", "gpt-4o", "Model to use for text generation")
}

func Execute(version string, commit string, date string, builtBy string) {
	rootCmd.Version = fmt.Sprintf("%s, commit %s, built at %s, build by %s", version, commit, date, builtBy)
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getAPIKey() (openai.APIKey, error) {
	apiKeyFilePath := os.Getenv("HOME") + "/.ai-text-shaper-apikey"
	bytes, err := os.ReadFile(apiKeyFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}
	return openai.APIKey(strings.TrimSuffix(string(bytes), "\n")), nil
}

func makeGenerativeAIClient(model string) (process.GenerativeAIClient, error) {
	apikey, err := getAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}
	return openai.New(apikey, model), nil
}

var rootCmd = &cobra.Command{
	Use:   "ai-text-shaper",
	Short: "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model",
	Long:  "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.New(&c, tui.Confirm).Run(args, makeGenerativeAIClient)
	},
}
