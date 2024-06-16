package runner

import (
	"bufio"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"os"
	"regexp"
	"strings"
)

func writeToFile(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func findMarkdownFirstCodeBlock(text string) (string, error) {
	re, err := regexp.Compile("(?s)```[a-zA-Z0-9]*?\n(.*?)\n```")
	if err != nil {
		return "", fmt.Errorf("error compiling regex: %w", err)
	}
	match := re.FindStringSubmatch(text)
	if match != nil {
		return match[1], nil
	}
	return "", nil
}

func diff(leftText, rightText string) string {
	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(leftText, rightText)
	diffs := dmp.DiffMain(a, b, false)
	diffs = dmp.DiffCharsToLines(diffs, c)
	return dmp.DiffPrettyText(diffs)
}

func confirm(s string, tries int) bool {
	r := bufio.NewReader(os.Stdin)

	for ; tries > 0; tries-- {
		fmt.Printf("%s [y/n]: ", s)

		res, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		// Empty input (i.e. "\n")
		if len(res) < 2 {
			continue
		}

		return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
	}

	return false
}

func (r *Runner) outputResult(resultText, inputText, outpath string) error {
	outputText := resultText
	if r.config.UseFirstCodeBlock {
		codeBlock, err := findMarkdownFirstCodeBlock(resultText)
		if err != nil {
			return fmt.Errorf("error finding first code block: %w", err)
		}
		if codeBlock != "" {
			outputText = codeBlock
		}
	}
	outputText = strings.TrimSuffix(outputText, "\n")

	if !r.config.Silent {
		fmt.Println(outputText)
		if r.config.Diff {
			fmt.Println(diff(inputText, outputText))
		} else {
		}
	}
	if outpath != "" {
		r.verboseLog("Writing to file: %s", outpath)
		if err := writeToFile(outpath, outputText); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}
	return nil
}
