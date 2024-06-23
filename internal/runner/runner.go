package runner

import (
	"errors" // fixed: Added for defining static errors
	"fmt"
	"github.com/ytka/ai-text-shaper/internal/ioutil"
	"github.com/ytka/ai-text-shaper/internal/steps"
	"log"
	"os"
)

var (
	ErrPromptOrPromptPathRequired = errors.New("either prompt or prompt-path must be provided")
	ErrOutpathRewriteConflict     = errors.New("outpath and rewrite cannot be provided together")
	ErrOutpathMultipleFiles       = errors.New("outpath cannot be provided when multiple input files are provided")
)

func (c *Config) Validate(inputFiles []string) error {
	if c.Prompt == "" && c.PromptPath == "" {
		return ErrPromptOrPromptPathRequired
	}
	if c.Outpath != "" && c.Rewrite {
		return ErrOutpathRewriteConflict
	}
	if c.Outpath != "" && len(inputFiles) > 1 {
		return ErrOutpathMultipleFiles
	}
	return nil
}

// Runner manages the execution of text processing tasks.
type Runner struct {
	config                         *Config
	inputFiles                     []string
	generativeAIHandlerFactoryFunc GenerativeAIHandlerFactoryFunc
	confirmFunc                    ConfirmFunc
}

type (
	GenerativeAIHandlerFactoryFunc func(model string) (steps.GenerativeAIClient, error)
	ConfirmFunc                    func(string) (bool, error)
)

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

type RunOption struct {
	gaiClient      steps.GenerativeAIClient
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
	promptText, err := steps.GetPromptText(r.config.Prompt, r.config.PromptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt text: %w", err)
	}
	pipeAvailable, err := ioutil.IsAvailablePipe(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to check if stdin is pipe: %w", err)
	}
	if pipeAvailable && len(r.inputFiles) >= 1 {
		added, err := steps.GetInputText("-")
		if err != nil {
			return nil, fmt.Errorf("failed to get input text from stdin: %w", err)
		}
		promptText += "\n" + added
	}

	var inputFilePaths []string
	if len(r.inputFiles) == 0 {
		inputFilePaths = []string{"-"}
	} else {
		inputFilePaths = r.inputFiles
	}

	return &RunOption{gaiClient: gai, promptText: promptText, inputFilePaths: inputFilePaths}, nil
}

// Run processing of multiple input files.
func (r *Runner) Run(opt *RunOption, onBeforeProcessing func(string), onAfterProcessing func(string)) error {
	for i, inputPath := range opt.inputFilePaths {
		p := NewProcess(r.config, r.confirmFunc)
		if err := p.Run(i, inputPath, opt, onBeforeProcessing, onAfterProcessing); err != nil {
			return fmt.Errorf("processing error: %w", err)
		}
	}
	return nil
}
