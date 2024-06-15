# AITextShaper

AITextShaper is a tool designed to shape and transform text using OpenAI's GPT-4 model. It provides functionalities to rewrite, generate, and print text based on given prompts and input files.

## Features

- **Rewrite Text**: Modify the content of input text based on a specific prompt.
- **Generate Text**: Create new text content using OpenAI's GPT-4.
- **Print Text**: Display the generated or rewritten text directly.

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/textcrafter.git
    cd textcrafter
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Set up your OpenAI API key:
    - Create a file named `.openai-apikey` in your home directory and paste your OpenAI API key into this file.

## Usage

```sh
Usage: textcrafter <command> <prompt file> <input file>
```

### Commands

- `rewrite`: Rewrite the input text based on the prompt.
- `gen`: Generate new text based on the prompt.
- `print`: Print the generated or rewritten text.

### Examples

1. **Rewrite text**:
    ```sh
    textcrafter rewrite prompts/rewrite-prompt.txt inputs/input.txt
    ```

2. **Generate text**:
    ```sh
    textcrafter gen prompts/generate-prompt.txt inputs/input.txt
    ```

3. **Print text**:
    ```sh
    textcrafter print prompts/print-prompt.txt inputs/input.txt
    ```

## Project Structure

```plaintext
textcrafter/
├── go.mod
├── go.sum
├── ai-client.go
├── main.go
├── examples/
│   ├── hello.go
│   ├── junit-test.kt
│   ├── hello-correct.go
│   ├── jackson.kt
│   ├── jackson-to-kseriazaltion.kt
├── prompts/
│   ├── correct.txt
│   ├── to-kotest.txt
│   ├── to-c.txt
│   ├── to-kseriazaltion.txt
└── .idea/
    ├── ai-text-shaper.iml
    ├── workspace.xml
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or new features.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- Thanks to OpenAI for providing the GPT-4 model.
- Inspired by various text transformation tools and the community's need for flexible text shaping.
