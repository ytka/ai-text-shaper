package runner

import "fmt"

type Config struct {
	Prompt     string
	PromptPath string

	Silent  bool
	Verbose bool
	Diff    bool

	Rewrite              bool
	Outpath              string
	UseFirstCodeBlock    bool
	ConfirmBeforeWriting bool
}

func (c *Config) Validate() error {
	if c.Prompt == "" && c.PromptPath == "" {
		return fmt.Errorf("either prompt or prompt-path must be provided")
	}
	if c.Outpath == "" && c.Rewrite {
		return fmt.Errorf("either outpath must be provided or rewrite must be true")
	}
	return nil
}
