package runner

import (
	"fmt"
)

type Config struct {
	Prompt                string
	PromptPath            string
	PromptOptimize        bool
	Model                 string
	MaxTokens             int
	MaxCompletionRepeatCount int
	DryRun                bool
	Silent                bool
	Verbose               bool
	Diff                  bool
	InputFileList         string
	LogAPILevel           string
	Rewrite               bool
	Outpath               string
	UseFirstCodeBlock     bool
	Confirm               bool
}

func (c *Config) Validate(inputFiles []string) error {
	if c.Prompt == "" && c.PromptPath == "" {
		return fmt.Errorf("either prompt or prompt-path must be provided")
	}
	if c.Outpath != "" && c.Rewrite {
		return fmt.Errorf("outpath and rewrite cannot be provided together")
	}
	if c.Outpath != "" && len(inputFiles) > 1 {
		return fmt.Errorf("outpath cannot be provided when multiple input files are provided")
	}
	return nil
}
