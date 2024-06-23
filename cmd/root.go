package cmd

import (
	"fmt"
	"github/ytka/ai-text-shaper/internal/steps"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"github/ytka/ai-text-shaper/internal/openai"
	"github/ytka/ai-text-shaper/internal/runner"
	"github/ytka/ai-text-shaper/internal/tui"
)

var (
	c       runner.Config
	rootCmd = &cobra.Command{
		Use:   "ai-text-shaper",
		Short: "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model",
		Long:  "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model.",
		RunE: func(cmd *cobra.Command, args []string) error {
			inputFiles := args
			if c.InputFileList != "" {
				files, err := readInputFiles(c.InputFileList)
				if err != nil {
					return err
				}
				inputFiles = files
			}
			return doRun(inputFiles, makeGAIFunc)
		},
	}
)

func init() {
	rootCmd.Version = "testX"

	// Prompt options
	rootCmd.Flags().StringVarP(&c.Prompt, "prompt", "p", "", "Prompt text")
	rootCmd.Flags().StringVarP(&c.PromptPath, "prompt-path", "P", "", "Prompt file path")
	rootCmd.Flags().BoolVarP(&c.PromptOptimize, "prompt-optimize", "O", true, "Optimize prompt text")

	// Model options
	rootCmd.Flags().StringVarP(&c.Model, "model", "m", "gpt-4o", "Model to use for text generation")
	rootCmd.Flags().IntVarP(&c.MaxTokens, "max-tokens", "t", 0, "Max tokens to generate")
	rootCmd.Flags().IntVar(&c.MaxCompletionRepeatCount, "max-completion-repeat-count", 1, "Max completion repeat count")

	// Stdout messages options
	rootCmd.Flags().BoolVarP(&c.DryRun, "dry-run", "D", false, "Dry run")
	rootCmd.Flags().BoolVarP(&c.Verbose, "verbose", "v", false, "Verbose output")
	rootCmd.Flags().BoolVarP(&c.Silent, "silent", "s", false, "Suppress output")
	rootCmd.Flags().BoolVarP(&c.Diff, "diff", "d", false, "Show diff of the input and output text")

	// Input file options
	rootCmd.Flags().StringVarP(&c.InputFileList, "input-file-list", "i", "", "Input file list")

	// Debug options
	rootCmd.Flags().StringVarP(&c.LogAPILevel, "log-api-level", "l", "", "API log level: info, debug")

	// Write file options
	rootCmd.Flags().BoolVarP(&c.Rewrite, "rewrite", "r", false, "Rewrite the input file with the result")
	rootCmd.Flags().StringVarP(&c.Outpath, "outpath", "o", "", "Output file path")
	rootCmd.Flags().BoolVarP(&c.UseFirstCodeBlock, "use-first-code-block", "f", false, "Use the first code block in the output text")
	rootCmd.Flags().BoolVarP(&c.Confirm, "confirm", "c", false, "Confirm before writing to file")
}

func Execute(version, commit, date, builtBy string) {
	var sb strings.Builder
	sb.WriteString(version)
	sb.WriteString(", commit ")
	sb.WriteString(commit)
	sb.WriteString(", built at ")
	sb.WriteString(date)
	sb.WriteString(", build by ")
	sb.WriteString(builtBy)
	rootCmd.Version = sb.String()
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

func makeGAIFunc(model string) (steps.GenerativeAIClient, error) {
	apikey, err := getAPIKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}
	var maxTokens *int
	if c.MaxTokens > 0 {
		maxTokens = &c.MaxTokens
	}
	return openai.New(apikey, model, c.LogAPILevel, maxTokens), nil
}

func readInputFiles(fileName string) ([]string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read input files from %s: %w", fileName, err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	files := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		files = append(files, line)
	}
	return files, nil
}

func isPipe(file *os.File) (bool, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return false, fmt.Errorf("failed to get file info: %w", err)
	}
	// Checks if the mode is pipe
	return (fileInfo.Mode() & os.ModeNamedPipe) != 0, nil
}

func doRun(inputFiles []string, makeGAIFunc func(model string) (steps.GenerativeAIClient, error)) error {
	r := runner.New(&c, inputFiles, makeGAIFunc, tui.Confirm)
	ropt, err := r.Setup()
	if err != nil {
		return fmt.Errorf("failed to setup runner: %w", err)
	}

	onBeforeProcessing := func(string) {}
	onAfterProcessing := func(string) {}

	pipe, err := isPipe(os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to check if stdin is pipe: %w", err)
	}

	if !pipe {
		var wg sync.WaitGroup
		var statusUI *tui.StatusUI
		onBeforeProcessing = func(inpath string) {
			wg.Add(1)
			input := "Stdin"
			if inpath != "-" {
				input = inpath
			}
			statusUI = tui.NewStatusUI(fmt.Sprintf("Processing... [%s]", input))
			go func() {
				defer wg.Done()
				if err := statusUI.Run(); err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "failed to run status UI: %v\n", err)
				}
			}()
		}
		onAfterProcessing = func(string) {
			statusUI.Quit()
			statusUI = nil
			wg.Wait()
		}
	}

	if err := r.Run(ropt, onBeforeProcessing, onAfterProcessing); err != nil {
		return fmt.Errorf("failed to run: %w", err)
	}
	return nil
}
