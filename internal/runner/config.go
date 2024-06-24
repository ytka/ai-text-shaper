package runner

type Config struct {
	Prompt                   string
	PromptPath               string
	PromptOptimize           bool
	Model                    string
	MaxTokens                int
	MaxCompletionRepeatCount int
	DryRun                   bool
	Silent                   bool
	Verbose                  bool
	ShowCost                 bool
	Diff                     bool
	InputFileList            string
	LogAPILevel              string
	Rewrite                  bool
	Outpath                  string
	UseFirstCodeBlock        bool
	Confirm                  bool
}

// Validate checks the configuration for errors.
func (c *Config) Validate(inputFiles []string) error {
	if c.Prompt == "" && c.PromptPath == "" {
		return ErrPromptOrPromptPathRequired
	}
	if c.Outpath != "" && c.Rewrite {
		return ErrOutpathRewriteConflict
	}
	if c.Outpath != "" && len(inputFiles) > 1 {
		return ErrOutpathMultipleFiles
	}
	return nil
}
