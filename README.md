# textforge

## Overview

`textforge` is a CLI tool for processing text files such as source code and Markdown using the OpenAI API.
How to process the text is instructed by providing an arbitrary text as a prompt, such as "translate to English".
By customizing the prompt, it can be used for various purposes such as code refactoring, translation, and conversion to a specified format.
Examples of prompts can be found in the files under the prompts/ directory.

`textforge` processes the given prompt into a format that is easy to handle and then sends it to OpenAI API's chat API, formats the result, and outputs it to the standard output.
In addition to outputting, you can also rewrite the original file, write it to a specified path, or confirm the output result before writing.

## Installation

To install `textforge`, download the binary from the following link:

[https://github.com/ytka/textforge/releases](https://github.com/ytka/textforge/releases)

## Set OpenAI API Key
This tool requires an API key for utilizing OpenAI API.
Create a file named `.textforge-apikey` in your home directory and write the OpenAI API key in it.

## Usage

The general usage pattern for `textforge` is as follows:

```sh
textforge [options] [input files...]
```

You can specify one or multiple input files. If no input files are specified, it reads from the standard input.

### Options

#### Prompt Options

- `-p, --prompt string`
   - Specify the prompt text.

- `-P, --prompt-path string`
   - Specify the path to the prompt file (text file). The string read from this file will be used as the prompt.

- `-m, --model string`
   - Specify the chat model to use. The default is `gpt-4o`.

#### Output Options

- `-v, --verbose`
   - Enable verbose output.

- `-s, --silent`
   - Suppress all output.

- `-d, --diff`
   - Display the difference between the input and output text along with the normal output.

- `-l, --log-api-level string`
   - Specify the API log level. `info` or `debug`.

- `-C, --show-cost`
   - Display the cost of text generation.

#### File Writing Options

- `-r, --rewrite`
   - Rewrite the input file with the result.

- `-o, --outpath string`
   - Specify the path of the output file.

- `-f, --use-first-code-block`
   - If the output text contains code blocks, use the first code block as the output.

- `-c, --confirm`
   - Ask for confirmation before writing to a file.

#### Other Options

- `-D, --dry-run`
   - Test operation without making any actual changes.

- `-t, --max-tokens int`
   - Specify the maximum number of tokens to generate.

- `--max-completion-repeat-count int`
   - Specify the maximum number of completion repeats (default 1).

- `-O, --prompt-optimize`
   - Optimize the prompt text (default true).

- `--version`
   - Display the version information of `textforge`.

## Examples

### Basic Usage

To specify a prompt from the command line:
```sh
textforge -p "prompt text" /path/to/inputfile.txt
```

### Use a Prompt File

To specify a prompt from a file:
```sh
textforge -P /path/to/promptfile.txt /path/to/inputfile.txt
```

### Verbose Output

To enable verbose output:

```sh
textforge -v /path/to/inputfile.txt
```

### Suppressing Output

To suppress all output:

```sh
textforge -s /path/to/inputfile.txt
```

### Display Differences

To display differences between the input and output text:

```sh
textforge -d /path/to/inputfile.txt
```

### Writing to a File

To write the result to a specific output file:

```sh
textforge -o /path/to/outputfile.txt /path/to/inputfile.txt
```

### Rewriting the Input File

To rewrite the input file with the result:

```sh
textforge -r /path/to/inputfile.txt
```

### Using the First Code Block

To use the first code block of the output text:

```sh
textforge -f /path/to/inputfile.txt
```

### Confirm Before Writing

To ask for confirmation before writing to a file:

```sh
textforge -c /path/to/inputfile.txt
```

## Examples of Use in Actual Development

In this project, textforge is used for development.
It is registered as a task in Taskfile.yaml, and tasks such as AI code review correction and translation processing are automated.
We will introduce the contents of the tasks to show what kind of automation can be done.

### lint-fix-ai

Automatically fix static analysis errors using AI. It can fix issues that cannot be fixed with the normal lint fix option (though not always).
```sh
task lint-fix-ai
```

### auto-review-changed

Review the changed files (git diff) and make corrections if necessary.
```sh
task auto-review-changed
```

### auto-review

Review the entire project and make corrections if necessary.
```sh
task auto-review
```

### auto-commit

Create a commit message based on the changes and commit. Confirm the message content before committing.
```sh
task auto-commit
```

### forge-doc-by-help

Update the README based on the program's help.
```sh
task forge-doc-by-help
```

### forge-doc-by-taskfile

Update the usage examples in the README based on the Taskfile usage examples.
```sh
task forge-doc-by-taskfile
```

### translate-prompts-to-en

Translate Japanese sample prompts into English.
```sh
task translate-prompts-to-en
```

### translate-readme-to-en

Translate the Japanese README into English.
```sh
task translate-readme-to-en
```

## License

This project is licensed under the [MIT License](link_to_license).
