package process

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func confirm(message string) (bool, error) {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s [y/N]: ", message)

		res, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return false, nil
			}
			return false, err
		}
		// Empty input (i.e. "\n")
		if len(res) < 2 {
			return false, nil
		}

		if strings.ToLower(strings.TrimSpace(res))[0] == 'y' {
			return true, nil
		}
	}
}

func WriteResult(outputText, outpath string, needConfirm bool) error {
	if needConfirm {
		conf, err := confirm("May I write the results of this operation to a file?")
		if err != nil {
			return err
		}
		if !conf {
			return nil
		}
	}
	if err := os.WriteFile(outpath, []byte(outputText), 0644); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}
