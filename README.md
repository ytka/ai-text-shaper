# textforge

## Overview

`textforge` is a CLI tool for processing text files such as source code and Markdown using the OpenAI API. How to process the text is specified by giving a prompt, such as "translate to English." By customizing the prompt, the tool can be used for various purposes such as code refactoring, translation, or converting to a specified format. Examples of prompts can be found in the files under the `prompts/` directory.

`textforge` formats the given prompt to make it easier to process before sending it to OpenAI's chat API, and then formats the result and outputs it to the standard output. Additionally, the tool can rewrite the original file, write to a specified path, or allow you to verify the output before writing.

## Installation

To install `textforge`, download the binaries from the following link:

[https://github.com/ytka/textforge/releases](https://github.com/ytka/textforge/releases)

## Setting the OpenAI API Key

This tool requires an API key to use the OpenAI API. Create a file named `.textforge-apikey` in your home directory and write your OpenAI API key in it.

## Usage

The general usage pattern for `textforge` is as follows:

```sh
textforge [options] [input files...]
```

You can specify one or more input files. If no input files are specified, it reads from the standard input.

### Options

#### General Options

- `-h, --help`
   - Show help for `textforge`.
- `--version`
   - Show version information for `textforge`.

#### Prompt Options

- `-p, --prompt string`
   - Specify the text for the prompt.
- `-P, --prompt-path string`
   - Specify the path to a prompt file (text file). The string read from this file will be used as the prompt.
- `-m, --model string`
   - Specify the model to use for Chat. The default is `gpt-4o`.
- `-t, --max-tokens int`
   - Specify the maximum number of tokens to generate.
- `-O, --prompt-optimize`
   - Optimize prompt text (default true).

#### Output Options

- `-v, --verbose`
   - Enable detailed output.
- `-s, --silent`
   - Suppress all output.
- `-d, --diff`
   - In addition to normal output, show the differences between the input and output texts.
- `-C, --show-cost`
   - Show the cost of the text generation.

#### Execution Options

- `-D, --dry-run`
   - Run without making actual changes.
- `-l, --log-api-level string`
   - Set API log level: info, debug.
- `--max-completion-repeat-count int`
   - Set the maximum number of completion repeats (default 1).

#### File Writing Options

- `-i, --input-file-list string`
   - Provide a list of input files.
- `-r, --rewrite`
   - Rewrite the input file with the result.
- `-o, --outpath string`
   - Specify the path to the output file.
- `-f, --use-first-code-block`
   - If the output text contains code blocks, use the first code block as the output.
- `-c, --confirm`
   - Ask for confirmation before writing to a file.

## Examples

### Basic Usage

To give a prompt from the command line:

```sh
textforge -p "prompt text" /path/to/inputfile.txt
```

### Using a Prompt File

To give a prompt from a file:

```sh
textforge -P /path/to/promptfile.txt /path/to/inputfile.txt
```

### Enabling Detailed Output

To enable detailed output:

```sh
textforge -v /path/to/inputfile.txt
```

### Suppressing Output

To suppress all output:

```sh
textforge -s /path/to/inputfile.txt
```

### Showing Differences

To show the differences between the input and output texts:

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

### Confirmation Before Writing

To ask for confirmation before writing to a file:

```sh
textforge -c /path/to/inputfile.txt
```

## License

This project is licensed under the [MIT License](link_to_license).
