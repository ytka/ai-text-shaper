package cmd

import "github.com/spf13/cobra"

type flags struct {
	prompt     string
	promptPath string

	verbose bool
	silent  bool
	diff    bool

	rewrite bool
	outpath string

	useFirstCodeBlock bool
}

func (f *flags) initCommandFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.prompt, "prompt", "p", "", "Prompt text")
	cmd.Flags().StringVarP(&f.promptPath, "prompt-path", "P", "", "Prompt file path")

	cmd.Flags().BoolVarP(&f.verbose, "verbose", "v", false, "Verbose mode")
	cmd.Flags().BoolVarP(&f.silent, "silent", "s", false, "Silent mode")
	cmd.Flags().BoolVarP(&f.diff, "diff", "d", false, "Show diff")

	cmd.Flags().BoolVarP(&f.rewrite, "rewrite", "r", false, "Rewrite the input file with the result")
	cmd.Flags().StringVarP(&f.outpath, "outpath", "o", "", "Output file path")

	cmd.Flags().BoolVarP(&f.useFirstCodeBlock, "use-first-code-block", "f", false, "Use the first code block in the output text")
}
