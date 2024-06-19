package runner

import (
	"ai-text-shaper/internal/process"
	"fmt"
	"log"
	"os"
	"strings"
)

// Runner manages the execution of text processing tasks.
type Runner struct {
	config                         *Config
	inputFiles                     []string
	generativeAIHandlerFactoryFunc GenerativeAIHandlerFactoryFunc
	confirmFunc                    ConfirmFUnc
}

type GenerativeAIHandlerFactoryFunc func(model string) (process.GenerativeAIClient, error)
type ConfirmFUnc func(string) (bool, error)

func New(config *Config,
	inputFiles []string,
	gaiFactory GenerativeAIHandlerFactoryFunc,
	confirmFunc ConfirmFUnc,
) *Runner {
	return &Runner{config: config, inputFiles: inputFiles, generativeAIHandlerFactoryFunc: gaiFactory, confirmFunc: confirmFunc}
}

func (r *Runner) verboseLog(msg string, args ...interface{}) {
	if r.config.Verbose {
		log.Printf(msg, args...)
	}
}

// runSingleInput processes a single input file using the GenerativeAIClient.
func (r *Runner) runSingleInput(index int, inputFilePath string, promptText string, gai process.GenerativeAIClient,
	onBeforeProcessing func(), onAfterProcessing func()) error {
	onBeforeProcessing()

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
	s := process.NewShaper(gai, r.config.MaxCompletionRepeatCount, r.config.UseFirstCodeBlock)
	result, err := s.ShapeText(promptText, inputText)
	if err != nil {
		return err
	}
	processedPromptText, rawResult, resultText := result.Prompt, result.RawResult, result.Result
	if !strings.HasSuffix(resultText, "\n") {
		resultText += "\n"
	}
	r.verboseLog("[%d] mergedPromptText: size:%d, '%s'", index, len(processedPromptText), processedPromptText)
	r.verboseLog("[%d] rawResult: size:%d, '%s'", index, len(rawResult), rawResult)
	r.verboseLog("[%d] resultText: '%s'", index, resultText)
	onAfterProcessing()

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

type RunOption struct {
	gaiClient      process.GenerativeAIClient
	promptText     string
	inputFilePaths []string
}

func (r *Runner) Setup() (*RunOption, error) {
	r.verboseLog("configs: %+v", r.config)
	r.verboseLog("inputFiles: %+v", r.inputFiles)
	if err := r.config.Validate(r.inputFiles); err != nil {
		return nil, fmt.Errorf("invalid configuration: %+v, %w", r.config, err)
	}
	r.verboseLog("make generative ai client")
	gai, err := r.generativeAIHandlerFactoryFunc(r.config.Model)
	if err != nil {
		return nil, fmt.Errorf("failed to make generative ai client: %w", err)
	}
	r.verboseLog("get prompt")
	promptText, err := process.GetPromptText(r.config.Prompt, r.config.PromptPath)
	if err != nil {
		return nil, err
	}
	r.verboseLog("promptText: '%s'", promptText)

	var inputFilePaths []string
	if len(r.inputFiles) == 0 {
		inputFilePaths = []string{"-"}
	} else {
		inputFilePaths = r.inputFiles
	}

	return &RunOption{gaiClient: gai, promptText: promptText, inputFilePaths: inputFilePaths}, nil
}

// Run processing of multiple input files
func (r *Runner) Run(opt *RunOption, onBeforeProcessing func(), onAfterProcessing func()) error {

	wrappedOnBeforeProcessStatus := func() {
		r.verboseLog("start processing")
		onBeforeProcessing()
	}
	wrappedOnAfterProcessStatus := func() {
		r.verboseLog("end processing")
		onAfterProcessing()
	}

	for i, inputPath := range opt.inputFilePaths {
		if err := r.runSingleInput(i+1, inputPath, opt.promptText, opt.gaiClient, wrappedOnBeforeProcessStatus, wrappedOnAfterProcessStatus); err != nil {
			return err
		}
	}
	return nil
}
