# ai-text-shaper

## Overview

`ai-text-shaper` is a CLI tool that uses the OpenAI API to process text files such as source code or Markdown according to a specified prompt. By customizing the prompt, you can perform tasks such as code refactoring, translation, or conversion to a specified format.

You can provide any text as a prompt, such as "translate to English". `ai-text-shaper` processes the given prompt, sends it via the OpenAI API's chat API, formats the result, and outputs it to standard output. Additionally, it can rewrite the original file, allow you to review the output before writing, and other functionalities. For examples of other prompts, refer to the files in the prompts/ directory.

## Installation

To install `ai-text-shaper`, download the binary from the following link:

[https://github.com/ytka/ai-text-shaper/releases](https://github.com/ytka/ai-text-shaper/releases)

## Setting up the OpenAI API Key
This tool requires an API key to use the OpenAI API. Create a file named `.ai-text-shaper-apikey` in your home directory and write your OpenAI API key in it.

## Usage

The general usage pattern for `ai-text-shaper` is as follows:

```sh
ai-text-shaper [options] [input_files...]
```

You can specify one or multiple input files. If no input files are specified, it reads from standard input.

### Options

#### Prompt Options

- `-p, --prompt string`
   - Specify the text for the prompt.

- `-P, --prompt-path string`
   - Specify the path to a prompt file (text file). The string read from this file will be used as the prompt.

#### Output Options

- `-v, --verbose`
   - Enable verbose output.

- `-s, --silent`
   - Suppress all output.

- `-d, --diff`
   - Show the diff between input and output text in addition to the standard output.

#### File Writing Options

- `-r, --rewrite`
   - Rewrite the input file with the result.

- `-o, --outpath string`
   - Specify the path for the output file.

- `-f, --use-first-code-block`
   - If the output text contains code blocks, use the first code block for output.

- `-c, --confirm`
   - Ask for confirmation before writing to a file.

## Examples

### Basic Usage

To provide a prompt from the command line:
```sh
ai-text-shaper -p "prompt text" /path/to/inputfile.txt
```

### Using a Prompt File

To provide a prompt from a file:
```sh
ai-text-shaper -P /path/to/promptfile.txt /path/to/inputfile.txt
```

### Enable Verbose Output

To enable verbose output:

```sh
ai-text-shaper -v /path/to/inputfile.txt
```

### Suppress Output

To suppress all output:

```sh
ai-text-shaper -s /path/to/inputfile.txt
```

### Show Diff

To show the diff between the input and output text:

```sh
ai-text-shaper -d /path/to/inputfile.txt
```

### Writing to a File

To write the result to a specific output file:

```sh
ai-text-shaper -o /path/to/outputfile.txt /path/to/inputfile.txt
```

### Rewriting the Input File

To rewrite the input file with the result:

```sh
ai-text-shaper -r /path/to/inputfile.txt
```

### Using the First Code Block

To use the first code block in the output text:

```sh
ai-text-shaper -f /path/to/inputfile.txt
```

### Confirm Before Writing

To confirm before writing to a file:

```sh
ai-text-shaper -c /path/to/inputfile.txt
```

## License

This project is licensed under the [MIT License](link_to_license).
