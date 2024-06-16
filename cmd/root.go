package cmd

import (
	"ai-text-shaper/internal/runner"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var c runner.Config

func init() {
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
	rootCmd.Flags().BoolVarP(&c.UseFirstCodeBlock, "use-first-code-block", "f", true, "Use the first code block in the output text")
	rootCmd.Flags().BoolVarP(&c.ConfirmBeforeWriting, "confirm-before-write", "c", false, "Confirm before writing to file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ai-text-shaper",
	Short: "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model",
	Long:  "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runner.New(&c).Run(args)
	},
}
