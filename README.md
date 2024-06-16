# ai-text-shaper

## Overview

`ai-text-shaper` is a CLI tool for processing text files such as source code or Markdown using the OpenAI API. It processes the text based on given prompts like "translate to English" or any other custom text. By customizing the prompts, you can use it for various purposes such as code refactoring, translation, or conversion to a specific format. See the files under prompts/ for examples of prompts.

`ai-text-shaper` processes the given prompt into a format that the OpenAI API's chat API can handle, then submits it and formats the result for standard output. Besides outputting, it can also rewrite the original file, write to a specified path, or allow you to review the output before writing.

## Installation

To install `ai-text-shaper`, download the binary from the following link:

[https://github.com/ytka/ai-text-shaper/releases](https://github.com/ytka/ai-text-shaper/releases)

## Set the OpenAI API Key

This tool requires an API key to use the OpenAI API. Create a file named `.ai-text-shaper-apikey` in your home directory and write the OpenAI API key in it.

## Usage

A common usage pattern for `ai-text-shaper` is as follows:

```sh
ai-text-shaper [options] [input file...]
```

You can specify one or multiple input files. If no input file is specified, it reads from standard input.

### Options

#### Prompt Options

- `-p, --prompt string`
   - Specify the text for the prompt.

- `-P, --prompt-path string`
   - Specify the path to a prompt file (text file). The string read from this file is used as the prompt.

#### Output Options

- `-v, --verbose`
   - Enable verbose output.

- `-s, --silent`
   - Suppress all output.

- `-d, --diff`
   - Show the difference between the input and output text in addition to the normal output.

#### File Writing Options

- `-r, --rewrite`
   - Rewrite the input file with the result.

- `-o, --outpath string`
   - Specify the path to the output file.

- `-f, --use-first-code-block`
   - Use the first code block in the output text.

- `-c, --confirm`
   - Ask for confirmation before writing to the file.

## Examples

### Basic Usage

To give a prompt from the command line:

```sh
ai-text-shaper -p "prompt text" /path/to/inputfile.txt
```

### Using a Prompt File

To give a prompt from a file:

```sh
ai-text-shaper -P /path/to/promptfile.txt /path/to/inputfile.txt
```

### Verbose Output

To enable verbose output:

```sh
ai-text-shaper -v /path/to/inputfile.txt
```

### Suppress Output

To suppress all output:

```sh
ai-text-shaper -s /path/to/inputfile.txt
```

### Show Differences

To show the difference between the input and output text:

```sh
ai-text-shaper -d /path/to/inputfile.txt
```

### Write to a File

To write the result to a specific output file:

```sh
ai-text-shaper -o /path/to/outputfile.txt /path/to/inputfile.txt
```

### Rewrite Input File

To rewrite the input file with the result:

```sh
ai-text-shaper -r /path/to/inputfile.txt
```

### Use First Code Block

To use the first code block in the output text:

```sh
ai-text-shaper -f /path/to/inputfile.txt
```

### Confirm Before Writing

To ask for confirmation before writing to the file:

```sh
ai-text-shaper -c /path/to/inputfile.txt
```

## License

This project is licensed under the [MIT License](link_to_license).
