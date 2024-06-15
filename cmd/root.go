package cmd

import (
	"ai-text-shaper/internal/iostore"
	"ai-text-shaper/internal/textshaper"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var fl flags

func init() {
	fl.initCommandFlags(rootCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	st := iostore.New(verboseLog)

	apikey, err := st.GetAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	inputFilePath := "-"
	if len(args) >= 1 {
		// FIXME: larger case
		inputFilePath = args[0]
	}

	outpath := fl.outpath
	if fl.rewrite {
		outpath = inputFilePath
	}

	promptText, err := st.GetPromptText(fl.prompt, fl.promptPath)
	if err != nil {
		return err
	}

	inputText, err := st.GetInputText(inputFilePath)
	if err != nil {
		return err
	}

	verboseLog("start shaping text")
	resultText, err := textshaper.ShapeText(apikey, promptText, inputText)
	verboseLog("end shaping text")
	if err != nil {
		return err
	}

	outputText := resultText
	if fl.useFirstCodeBlock {
		codeBlock, err := iostore.FindMarkdownFirstCodeBlock(resultText)
		if err != nil {
			return fmt.Errorf("error finding first code block: %w", err)
		}
		if codeBlock != "" {
			outputText = codeBlock
		}
	}
	outputText = strings.TrimSuffix(outputText, "\n")

	if !fl.silent {
		fmt.Println(outputText)
		if fl.diff {
			fmt.Println(iostore.Diff(inputText, outputText))
		} else {
		}
	}
	if outpath != "" {
		if err := st.WriteToFile(outpath, outputText); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "ai-text-shaper",
	Short: "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model",
	Long:  "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model.",
	RunE: func(cmd *cobra.Command, args []string) error {
		verboseLog("ai-text-shaper started")
		verboseLog("flags: %+v", fl)
		verboseLog("args: %v", args)
		return run(args)
	},
}
