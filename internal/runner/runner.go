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
	inputFiles                     []string
	generativeAIHandlerFactoryFunc GenerativeAIHandlerFactoryFunc
	confirmFunc                    ConfirmFunc
}

type GenerativeAIHandlerFactoryFunc func(model string) (process.GenerativeAIClient, error)
type ConfirmFunc func(string) (bool, error)

func New(config *Config, inputFiles []string, gaiFactory GenerativeAIHandlerFactoryFunc, confirmFunc ConfirmFunc) *Runner {
	return &Runner{
		config:                         config,
		inputFiles:                     inputFiles,
		generativeAIHandlerFactoryFunc: gaiFactory,
		confirmFunc:                    confirmFunc,
	}
}

func (r *Runner) verboseLog(msg string, args ...interface{}) {
	if r.config.Verbose {
		log.Printf(msg, args...)
	}
}

func (r *Runner) process(index int, inputFilePath string, promptText string, gai process.GenerativeAIClient) (*process.ShapeResult, error) {
	r.verboseLog("\n")
	r.verboseLog("[%d] get input text from: %s", index, inputFilePath)
	inputText, err := process.GetInputText(inputFilePath)
	if err != nil {
		return nil, err
	}
	r.verboseLog("[%d] inputText: '%s'", index, inputText)

	r.verboseLog("[%d] shaping text", index)
	shaper := process.NewShaper(gai, r.config.MaxCompletionRepeatCount, r.config.UseFirstCodeBlock, r.config.PromptOptimize)
	result, err := shaper.ShapeText(promptText, inputText)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Runner) output(shapeResult *process.ShapeResult, index int, inputFilePath string, inputText string) error {
	r.verboseLog("[%d] mergedPromptText: size:%d, '%s'", index, len(shapeResult.Prompt), shapeResult.Prompt)
	r.verboseLog("[%d] rawResult: size:%d, '%s'", index, len(shapeResult.RawResult), shapeResult.RawResult)
	r.verboseLog("[%d] resultText: '%s'", index, shapeResult.Result)

	if r.config.Rewrite {
		fmt.Println("Rewrite file:", inputFilePath)
	} else {
		if !r.config.Silent && !r.config.DryRun {
			process.OutputToStdout(shapeResult.Result, inputText, r.config.Diff)
		}
	}

	if r.config.DryRun {
		return nil
	}

	outpath := r.config.Outpath
	if r.config.Rewrite {
		outpath = inputFilePath
	}

	if r.config.Confirm {
		r.verboseLog("[%d] Confirming", index)
		conf, err := r.confirmFunc("Continue (y/n)?: ")
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
		return process.WriteResult(shapeResult.Result, outpath)
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
func (r *Runner) Run(opt *RunOption, onBeforeProcessing func(string), onAfterProcessing func(string)) error {
	for i, inputPath := range opt.inputFilePaths {
		r.verboseLog("start processing")

		onBeforeProcessing(inputPath)
		shapeResult := &process.ShapeResult{}
		if !r.config.DryRun {
			result, err := r.process(i+1, inputPath, opt.promptText, opt.gaiClient)
			r.verboseLog("end processing")
			if err != nil {
				onAfterProcessing(inputPath)
				return err
			}
			shapeResult = result
		}
		onAfterProcessing(inputPath)

		if err := r.output(shapeResult, i+1, inputPath, shapeResult.Prompt); err != nil {
			return err
		}
	}
	return nil
}
