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
	fmt.Print(outputText)

	if useDiff {
		fmt.Printf("\n====begin of diff==== in size: %d, out size: %d\n", len(inputText), len(outputText))
		fmt.Print(diff(inputText, outputText))
		fmt.Println("====end of diff====")
	}
}
