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
	"sync"
	"time"
)

var c runner.Config

func init() {
	rootCmd.Version = "testX"
	// prompt options
	rootCmd.Flags().StringVarP(&c.Prompt, "prompt", "p", "", "Prompt text")
	rootCmd.Flags().StringVarP(&c.PromptPath, "prompt-path", "P", "", "Prompt file path")

	// model options
	rootCmd.Flags().StringVarP(&c.Model, "model", "m", "gpt-4o", "statusModel to use for text generation")
	rootCmd.Flags().IntVarP(&c.MaxTokens, "max-tokens", "t", 0, "Max tokens to generate")
	rootCmd.Flags().IntVar(&c.MaxCompletionRepeatCount, "max-completion-repeat-count", 1, "Max completion repeat count")

	// stdout messages options
	rootCmd.Flags().BoolVarP(&c.Verbose, "verbose", "v", false, "Verbose output")
	rootCmd.Flags().BoolVarP(&c.Silent, "silent", "s", false, "Suppress output")
	rootCmd.Flags().BoolVarP(&c.Diff, "diff", "d", false, "Show diff of the input and output text")

	// input file options
	rootCmd.Flags().StringVarP(&c.InputFileList, "input-file-list", "i", "", "Input file list")

	// debug options
	rootCmd.Flags().StringVarP(&c.LogAPILevel, "log-api-level", "l", "", "API log level: info, debug")

	// write file options
	rootCmd.Flags().BoolVarP(&c.Rewrite, "rewrite", "r", false, "Rewrite the input file with the result")
	rootCmd.Flags().StringVarP(&c.Outpath, "outpath", "o", "", "Output file path")
	rootCmd.Flags().BoolVarP(&c.UseFirstCodeBlock, "use-first-code-block", "f", false, "Use the first code block in the output text")
	rootCmd.Flags().BoolVarP(&c.Confirm, "confirm", "c", false, "Confirm before writing to file")

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

func makeGAIFunc(model string) (process.GenerativeAIClient, error) {
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
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	files := make([]string, 0, len(lines))
	for _, line := range lines {
		line := strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		files = append(files, line)
	}
	return files, nil
}

var rootCmd = &cobra.Command{
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
		//return doRun(inputFiles, makeGAIFunc)
		return doRunWithStatus(inputFiles, makeGAIFunc)
	},
}

func doRun(args []string, makeGAIFunc func(model string) (process.GenerativeAIClient, error)) error {
	onChangeStatus := func(status string) {
		// fmt.Println(status)
	}

	return runner.New(&c, makeGAIFunc, tui.Confirm, onChangeStatus).
		Run(args)
}

func doRunWithStatus(args []string, makeGAIFunc func(model string) (process.GenerativeAIClient, error)) error {
	statusUI := tui.NewStatusUI("ai-text-shaper")

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := statusUI.Run(); err != nil {
			errChan <- err
		}
	}()

	statusUI.UpdateStatusText("start msg")
	onChangeStatus := func(status string) {
		statusUI.UpdateStatusText(status)
	}

	time.Sleep(1 * time.Second)
	onChangeStatus("after 1")
	time.Sleep(1 * time.Second)
	//time.Sleep(1 * time.Second)

	/*
		r := runner.New(&c, makeGAIFunc, tui.Confirm, onChangeStatus)
		if err := r.Run(args); err != nil {
			return err
		}
	*/

	statusUI.Quit()
	fmt.Println("hoge fuga")
	fmt.Println("dfafad")
	wg.Wait()

	return nil

	/*
		statusUI.Quit()

		wg.Wait()
		select {
		case runnerErr := <-errChan:
			return runnerErr
		default:
			return nil
		}

	*/
}
