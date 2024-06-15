# ai-text-shaper

ai-text-shaperは、OpenAIのGPTモデルを使用してテキストを整形および変換するツールです。指定されたプロンプトおよび入力ファイルに基づいて、テキストの書き換え、生成、および表示機能を提供します。

## 特長

- **テキストの書き換え**: 特定のプロンプトに基づいて入力テキストの内容を修正します。
- **テキストの生成**: OpenAIのGPT-4を使用して新しいテキストを作成します。
- **テキストの表示**: 生成されたまたは書き換えられたテキストを直接表示します。

## インストール

1. リポジトリをクローンします:
    ```sh
    git clone https://github.com/yourusername/textcrafter.git
    cd textcrafter
    ```

2. 依存関係をインストールします:
    ```sh
    go mod tidy
    ```

3. OpenAI APIキーを設定します:
    - 自分のホームディレクトリに`.openai-apikey`という名前のファイルを作成し、そのファイルにOpenAI APIキーを貼り付けます。

## 使用方法

```sh
Usage: textcrafter <command> <prompt file> <input file>
```

### コマンド

- `rewrite`: プロンプトに基づいて入力テキストを書き換えます。
- `gen`: プロンプトに基づいて新しいテキストを生成します。
- `print`: 生成されたまたは書き換えられたテキストを表示します。

### 例

1. **テキストの書き換え**:
    ```sh
    textcrafter rewrite prompts/rewrite-prompt.txt inputs/input.txt
    ```

2. **テキストの生成**:
    ```sh
    textcrafter gen prompts/generate-prompt.txt inputs/input.txt
    ```

3. **テキストの表示**:
    ```sh
    textcrafter print prompts/print-prompt.txt inputs/input.txt
    ```

## プロジェクト構造

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

## コントリビュート

貢献は歓迎されます！改善や新機能については、issueを開くかプルリクエストを提出してください。

## ライセンス

このプロジェクトはMITライセンスの下でライセンスされています。詳細については[LICENSE](LICENSE)ファイルをご覧ください。

## 謝辞

- GPT-4モデルを提供してくれたOpenAIに感謝します。
- 様々なテキスト変換ツールとコミュニティの柔軟なテキスト整形のニーズにインスパイアされました。