package runner

import (
	"fmt"
	"github/ytka/ai-text-shaper/internal/steps"
	"log"
	"os"
)

type Config struct {
	Prompt                   string
	PromptPath               string
	PromptOptimize           bool
	Model                    string
	MaxTokens                int
	MaxCompletionRepeatCount int
	DryRun                   bool
	Silent                   bool
	Verbose                  bool
	Diff                     bool
	InputFileList            string
	LogAPILevel              string
	Rewrite                  bool
	Outpath                  string
	UseFirstCodeBlock        bool
	Confirm                  bool
}

type Process struct {
	config      *Config
	confirmFunc ConfirmFunc
}

func NewProcess(config *Config, confirmFunc ConfirmFunc) *Process {
	return &Process{config: config, confirmFunc: confirmFunc}
}

func (p *Process) verboseLog(msg string, args ...interface{}) {
	if p.config.Verbose {
		log.Printf(msg, args...)
	}
}

func (p *Process) Run(i int, inputPath string, opt *RunOption, onBeforeProcessing func(string), onAfterProcessing func(string)) error {
	p.verboseLog("start processing")

	onBeforeProcessing(inputPath)
	shapeResult := &steps.ShapeResult{}
	if !p.config.DryRun {
		result, err := p.getInputAndShape(i+1, inputPath, opt.promptText, opt.gaiClient)
		p.verboseLog("end processing")
		if err != nil {
			onAfterProcessing(inputPath)
			return err
		}
		shapeResult = result
	}
	onAfterProcessing(inputPath)

	if err := p.output(shapeResult, i+1, inputPath, shapeResult.Prompt); err != nil {
		return err
	}
	return nil
}

func (p *Process) getInputAndShape(index int, inputFilePath string, promptText string, gai steps.GenerativeAIClient) (*steps.ShapeResult, error) {
	p.verboseLog("\n")
	p.verboseLog("[%d] get input text from: %s", index, inputFilePath)
	inputText, err := steps.GetInputText(inputFilePath)
	if err != nil {
		return nil, err
	}
	p.verboseLog("[%d] inputText: '%s'", index, inputText)

	p.verboseLog("[%d] shaping text", index)
	shaper := steps.NewShaper(gai, p.config.MaxCompletionRepeatCount, p.config.UseFirstCodeBlock, p.config.PromptOptimize)
	result, err := shaper.ShapeText(promptText, inputText)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *Process) confirm(index int, inputFilePath string) error {
	p.verboseLog("[%d] Confirming", index)
	conf, err := p.confirmFunc("Continue (y/n)?: ")
	if err != nil {
		return err
	}
	p.verboseLog("[%d] Confirmation: %t", index, conf)
	if !conf && inputFilePath == "-" {
		os.Exit(1)
	}
	return nil
}

func (p *Process) write(index int, resultText string, outpath string) error {
	if p.config.Rewrite {
		if p.config.DryRun {
			fmt.Printf("Rewrite file:%s, dry-run skipped.\n", outpath)
		} else {
			fmt.Printf("Rewrite file:%s\n", outpath)
		}
	}
	if outpath != "" && !p.config.DryRun {
		p.verboseLog("[%d] Writing to file: %s", index, outpath)
		if err := steps.WriteResult(resultText, outpath); err != nil {
			return err
		}
	}
	return nil
}

func (p *Process) output(shapeResult *steps.ShapeResult, index int, inputFilePath string, inputText string) error {
	p.verboseLog("[%d] mergedPromptText: size:%d, '%s'", index, len(shapeResult.Prompt), shapeResult.Prompt)
	p.verboseLog("[%d] rawResult: size:%d, '%s'", index, len(shapeResult.RawResult), shapeResult.RawResult)
	p.verboseLog("[%d] resultText: '%s'", index, shapeResult.Result)

	if !p.config.Silent && !p.config.DryRun && !p.config.Rewrite {
		steps.Print(shapeResult.Result, inputText, p.config.Diff)
	}

	if p.config.Confirm {
		if err := p.confirm(index, inputFilePath); err != nil {
			return err
		}
	}

	outpath := p.config.Outpath
	if p.config.Rewrite {
		outpath = inputFilePath
	}
	return p.write(index, shapeResult.Result, outpath)
}
