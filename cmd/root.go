package cmd

import (
	"ai-text-shaper/internal/textshaper"
	"ai-text-shaper/internal/tio"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"os"
)

var fl flags

func init() {
	fl.initCommandFlags(rootCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ai-text-shaper",
	Short: "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model",
	Long:  "ai-text-shaper is a tool designed to shape and transform text using OpenAI's GPT model.",
	// Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		verboseLog("ai-text-shaper started")
		verboseLog("flags: %+v", fl)
		verboseLog("args: %v", args)

		apikey, err := tio.GetAPIKey()
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

		verboseLog("start reading prompt")
		promptText, err := tio.GetPromptText(fl.prompt, fl.promptPath)
		if err != nil {
			return err
		}

		verboseLog("start reading input: %s", inputFilePath)
		inputText, err := tio.GetInputText(inputFilePath)
		if err != nil {
			return err
		}

		verboseLog("start shaping text")
		resultText, err := textshaper.ShapeText(apikey, promptText, inputText)
		verboseLog("end shaping text")
		if err != nil {
			return err
		}

		if !fl.silent {
			if fl.diff {
				dmp := diffmatchpatch.New()
				a, b, c := dmp.DiffLinesToChars(inputText, resultText)
				diffs := dmp.DiffMain(a, b, false)
				diffs = dmp.DiffCharsToLines(diffs, c)
				fmt.Println(dmp.DiffPrettyText(diffs))
			} else {
				fmt.Println(resultText)
			}
		}
		if outpath != "" {
			verboseLog("writing to file: %s", outpath)
			if err := tio.WriteToFile(outpath, resultText); err != nil {
				return fmt.Errorf("error writing to file: %w", err)
			}
		}

		return nil
	},
}
