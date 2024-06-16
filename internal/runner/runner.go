package runner

import (
	"ai-text-shaper/internal/textshaper"
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
	r.verboseLog("get API key")
	apikey, err := getAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}
	inputFilePath := "-"
	if len(inputFiles) >= 1 {
		// FIXME: larger case
		inputFilePath = inputFiles[0]
	}
	r.verboseLog("get prompt")
	promptText, err := getPromptText(r.config.Prompt, r.config.PromptPath)
	if err != nil {
		return err
	}
	r.verboseLog("get input")
	inputText, err := getInputText(inputFilePath)
	if err != nil {
		return err
	}

	/*
		Shape
	*/
	r.verboseLog("start shaping text")
	resultText, err := textshaper.ShapeText(string(apikey), promptText, inputText)
	r.verboseLog("end shaping text")
	if err != nil {
		return err
	}

	/*
		Output
	*/
	outpath := r.config.Outpath
	if r.config.Rewrite {
		outpath = inputFilePath
	}
	return r.outputResult(resultText, inputText, outpath)
}
