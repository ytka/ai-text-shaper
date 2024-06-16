package process

import (
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func diff(leftText, rightText string) string {
	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(leftText, rightText)
	diffs := dmp.DiffMain(a, b, false)
	diffs = dmp.DiffCharsToLines(diffs, c)
	return dmp.DiffPrettyText(diffs)
}

func OutputToStdout(outputText, inputText string, useDiff bool) {
	fmt.Println(outputText)
	if useDiff {
		fmt.Println(diff(inputText, outputText))
	}
}
