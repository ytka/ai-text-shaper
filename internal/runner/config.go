package runner

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
