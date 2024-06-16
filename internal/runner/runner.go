package runner

import (
	"ai-text-shaper/internal/openai"
	"ai-text-shaper/internal/process"
	"fmt"
	"log"
)

type Runner struct {
	config *Config
}

func New(config *Config) *Runner {
	return &Runner{config: config}
}

func (r *Runner) verboseLog(msg string, args ...interface{}) {
	if r.config.Verbose {
		log.Printf(msg, args...)
	}
}

func (r *Runner) Run(inputFiles []string) error {
	r.verboseLog("start run")
	r.verboseLog("configs: %+v", r.config)
	r.verboseLog("inputFiles: %+v", inputFiles)

	/*
		Prepare
	*/
	r.verboseLog("get OpenAI API key")
	apikey, err := openai.GetAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}
	inputFilePath := "-"
	if len(inputFiles) >= 1 {
		// FIXME: larger case
		inputFilePath = inputFiles[0]
	}
	r.verboseLog("get prompt")
	promptText, err := process.GetPromptText(r.config.Prompt, r.config.PromptPath)
	if err != nil {
		return err
	}
	r.verboseLog("get input")
	inputText, err := process.GetInputText(inputFilePath)
	if err != nil {
		return err
	}

	/*
		Shape
	*/
	r.verboseLog("start shaping text")
	resultText, err := process.ShapeText(string(apikey), promptText, inputText, r.config.UseFirstCodeBlock)
	r.verboseLog("end shaping text")
	if err != nil {
		return err
	}

	/*
		Output
	*/
	if !r.config.Silent {
		process.OutputToStdout(resultText, inputText, r.config.Diff)
	}
	outpath := r.config.Outpath
	if r.config.Rewrite {
		outpath = inputFilePath
	}
	if outpath != "" {
		r.verboseLog("Writing to file: %s", outpath)
		return process.WriteResult(resultText, outpath, r.config.ConfirmBeforeWriting)
	}
	return nil
}
