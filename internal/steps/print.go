package steps

import (
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// diff returns a formatted string showing the differences between leftText and rightText.
func diff(leftText, rightText string) string {
	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(leftText, rightText)
	diffs := dmp.DiffMain(a, b, false)
	diffs = dmp.DiffCharsToLines(diffs, c)
	return dmp.DiffPrettyText(diffs)
}

// GetDiffSize returns a boolean indicating if there are any differences and the number of characters added and removed.
func GetDiffSize(leftText, rightText string) (bool, int, int) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(leftText, rightText, false)

	added := 0
	removed := 0
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			added += len(diff.Text)
		case diffmatchpatch.DiffDelete:
			removed += len(diff.Text)
		case diffmatchpatch.DiffEqual:
			// No changes needed for equal parts
		}
	}

	return added > 0 || removed > 0, added, removed
}

// Print outputs the provided outputText to the console. If useDiff is true, it also prints the differences between inputText and outputText.
func Print(outputText, inputText string, useDiff bool) {
	fmt.Print(outputText)

	if useDiff {
		fmt.Printf(
			"\n====begin of diff==== in size: %d, out size: %d\n",
			len(inputText),
			len(outputText),
		)
		fmt.Print(diff(inputText, outputText))
		fmt.Println("====end of diff====")
	}
}
