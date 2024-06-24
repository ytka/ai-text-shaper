# textforge

## 概要

`textforge`は、ソースコードや Markdownなどのテキストファイルを OpenAI APIを使って加工するための CLIツールです。
どのように加工するかは、"英語に翻訳して"など任意のテキストをプロンプトとして与えることで指示します。
プロンプトをカスタマイズすることで、コードのリファクタリングや翻訳、指定した書式への変換などさまざまな用途に利用できます。
プロンプトの例は、prompts/以下のファイルを参照してください。

`textforge`は与えられたプロンプトを処理しやすい形に加工してから OpenAI APIの chat APIで送信し、結果を整形して標準出力に出力します。
出力の他にも、元のファイルを書き換える、指定したパスに書き出す、書き換え前に出力結果を確認してから書き出すこともできます。

## インストール

`textforge`をインストールするには、以下のリンクからバイナリをダウンロードしてください：

[https://github.com/ytka/textforge/releases](https://github.com/ytka/textforge/releases)

## OpenAI APIのAPIキーを設定
このツールはOpenAI APIを利用するため APIキーが必要です。
ホームディレクトリに `.textforge-apikey` という名前のファイルを作成し、OpenAI APIキーを書き込んでください。

## 使い方

`textforge`の一般的な使用パターンは以下の通りです：

```sh
textforge [オプション] [入力ファイル...]
```

入力ファイルは一つもしくは複数のファイルを指定できます。入力ファイルが指定されない場合、標準入力から読み取ります。

### オプション

#### プロンプトオプション

- `-p, --prompt string`
   - プロンプトのテキストを指定します。

- `-P, --prompt-path string`
   - プロンプトファイル（テキストファイル）のパスを指定します。このファイルから読み取った文字列をプロンプトとして使用します。

- `-m, --model string`
   - 使用するChat用モデルを指定します。デフォルトは `gpt-4o` です。

#### 出力オプション

- `-v, --verbose`
   - 詳細出力を有効にします。

- `-s, --silent`
   - すべての出力を抑制します。

- `-d, --diff`
   - 通常出力に加えて、入力と出力のテキストの差分を表示します。

- `-l, --log-api-level string`
   - APIログレベルを指定します。 `info` または `debug`。

- `-C, --show-cost`
   - テキスト生成のコストを表示します。

#### ファイル書き込みオプション

- `-r, --rewrite`
   - 結果で入力ファイルを書き換えます。

- `-o, --outpath string`
   - 出力ファイルのパスを指定します。

- `-f, --use-first-code-block`
   - 出力テキストにコードブロックが含まれる場合、最初のコードブロックを出力として使用します。

- `-c, --confirm`
   - ファイルに書き込む前に書き込んでよいか確認を求めます。

#### その他のオプション

- `-D, --dry-run`
   - 実際には変更を加えず動作をテストします。

- `-t, --max-tokens int`
   - 生成する最大トークン数を指定します。

- `--max-completion-repeat-count int`
   - 最大のコンプリート繰り返し回数を指定します（デフォルト 1）。

- `-O, --prompt-optimize`
   - プロンプトテキストの最適化を行います（デフォルト true）。

- `--version`
   - `textforge`のバージョン情報を表示します。

## 使用例

### 基本的な使用方法

プロンプトをコマンドラインから与えるには：
```sh
textforge -p "プロンプトのテキスト" /path/to/inputfile.txt
```

### プロンプトファイルを使用

プロンプトをファイルから与えるには：
```sh
textforge -P /path/to/promptfile.txt /path/to/inputfile.txt
```

### 詳細出力

詳細出力を有効にするには：

```sh
textforge -v /path/to/inputfile.txt
```

### 出力の抑制

全ての出力を抑制するには：

```sh
textforge -s /path/to/inputfile.txt
```

### 差分の表示

入力と出力のテキストの差分を表示するには：

```sh
textforge -d /path/to/inputfile.txt
```

### ファイルへの書き込み

結果を特定の出力ファイルに書き込むには：

```sh
textforge -o /path/to/outputfile.txt /path/to/inputfile.txt
```

### 入力ファイルの書き換え

入力ファイルを結果で書き換える(rewrite)には：

```sh
textforge -r /path/to/inputfile.txt
```

### 最初のコードブロックの使用

出力テキストの最初のコードブロックを使用するには：

```sh
textforge -f /path/to/inputfile.txt
```

### 書き込み前の確認

ファイルに書き込む前に確認するには：

```sh
textforge -c /path/to/inputfile.txt
```

## 実際の開発での利用例

このプロジェクトでは開発に textforge を使っています。
Taskfile.yamlにタスクとして登録し、タスクにはAIコードレビューによる修正や翻訳処理などを自動化しています。
どのような自動化ができるかタスクの内容を紹介します。

### lint-fix-ai

AIを使用して静的解析のエラーを自動修正します。通常の lint fixオプションでは修正できない問題も修正できます（できない場合もある）。
```sh
task lint-fix-ai
```

### auto-review-changed

変更されたファイル(git diff)をレビューし、必要があれば修正します。
```sh
task auto-review-changed
```

### auto-review

プロジェクト全体をレビューし、必要があれば修正します。
```sh
task auto-review
```

### auto-commit

変更を内容からコミットメッセージを作成し、コミットします。コミット前にメッセージ内容でよいか確認します。
```sh
task auto-commit
```

### forge-doc-by-help

プログラムのヘルプに基づいて、READMEを更新します。
```sh
task forge-doc-by-help
```

### forge-doc-by-taskfile

Taskfileの利用例に基づいて、READMEの使い方や利用例を更新します。
```sh
task forge-doc-by-taskfile
```

### translate-prompts-to-en

日本語のサンプルプロンプトを英語に翻訳します。
```sh
task translate-prompts-to-en
```

### translate-readme-to-en

日本語のREADMEを英語に翻訳します。
```sh
task translate-readme-to-en
```

## ライセンス

このプロジェクトは[MITライセンス](link_to_license)の下でライセンスされています。
