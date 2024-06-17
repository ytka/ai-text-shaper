package runner

import (
	"ai-text-shaper/internal/process"
	"fmt"
	"log"
	"os"
)

// Runner manages the execution of text processing tasks.
type Runner struct {
	config                         *Config
	generativeAIHandlerFactoryFunc GenerativeAIHandlerFactoryFunc
	confirmFunc                    ConfirmFUnc
}

type GenerativeAIHandlerFactoryFunc func(model string) (process.GenerativeAIClient, error)
type ConfirmFUnc func(string) (bool, error)

func New(config *Config, gaiFactory GenerativeAIHandlerFactoryFunc, confirmFunc ConfirmFUnc) *Runner {
	return &Runner{config: config, generativeAIHandlerFactoryFunc: gaiFactory, confirmFunc: confirmFunc}
}

func (r *Runner) verboseLog(msg string, args ...interface{}) {
	if r.config.Verbose {
		log.Printf(msg, args...)
	}
}

// runSingleInput processes a single input file using the GenerativeAIClient.
func (r *Runner) runSingleInput(index int, inputFilePath string, promptText string, gai process.GenerativeAIClient) error {
	r.verboseLog("\n")
	r.verboseLog("[%d] get input text from: %s", index, inputFilePath)
	inputText, err := process.GetInputText(inputFilePath)
	if err != nil {
		return err
	}
	r.verboseLog("[%d] inputText: '%s'", index, inputText)

	/*
		Shape
	*/
	r.verboseLog("[%d] shaping text", index)
	processedPromptText, rawResult, resultText, err := process.ShapeText(gai, promptText, inputText, r.config.UseFirstCodeBlock)
	if err != nil {
		return err
	}
	r.verboseLog("[%d] mergedPromptText: size:%d, '%s'", index, len(processedPromptText), processedPromptText)
	r.verboseLog("[%d] rawResult: size:%d, '%s'", index, len(rawResult), rawResult)
	r.verboseLog("[%d] resultText: '%s'", index, resultText)

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

	if r.config.Confirm {
		r.verboseLog("[%d] Confirming", index)
		conf, err := r.confirmFunc("Continue (y/N)?: ")
		if err != nil {
			return err
		}
		r.verboseLog("[%d] Confirmation: %t", index, conf)
		if !conf {
			os.Exit(1)
		}
	}
	if outpath != "" {
		r.verboseLog("[%d] Writing to file: %s", index, outpath)
		return process.WriteResult(resultText, outpath)
	}
	return nil
}

// Run processing of multiple input files
func (r *Runner) Run(inputFiles []string) error {
	r.verboseLog("start run")
	r.verboseLog("configs: %+v", r.config)
	r.verboseLog("inputFiles: %+v", inputFiles)

	if err := r.config.Validate(inputFiles); err != nil {
		return fmt.Errorf("invalid configuration: %+v, %w", r.config, err)
	}

	/*
		Prepare
	*/
	r.verboseLog("make generative ai client")
	gai, err := r.generativeAIHandlerFactoryFunc(r.config.Model)
	if err != nil {
		return fmt.Errorf("failed to make generative ai client: %w", err)
	}
	r.verboseLog("get prompt")
	promptText, err := process.GetPromptText(r.config.Prompt, r.config.PromptPath)
	if err != nil {
		return err
	}
	r.verboseLog("promptText: '%s'", promptText)

	/*
		Process
	*/
	var inputFilePaths []string
	if len(inputFiles) == 0 {
		inputFilePaths = []string{"-"}
	} else {
		inputFilePaths = inputFiles
	}
	for i, inputPath := range inputFilePaths {
		err := r.runSingleInput(i+1, inputPath, promptText, gai)
		if err != nil {
			return err
		}
	}

	return nil
}
