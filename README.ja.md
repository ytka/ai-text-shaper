# ai-text-shaper

## 概要

`ai-text-shaper`は、ソースコードや Markdownなどのテキストファイルを OpenAI APIを使って加工するための CLIツールです。
どのように加工するかは、"英語に翻訳して"など任意のテキストをプロンプトとして与えることで指示します。
プロンプトをカスタマイズすることで、コードのリファクタリングや翻訳、指定した書式への変換などさまざまな用途に利用できます。
プロンプトの例は、prompts/以下のファイルを参照してください。

`ai-text-shaper`は与えられたプロンプトを処理しやすい形に加工してから OpenAI APIの chat APIで送信し、結果を整形して標準出力に出力します。
出力の他にも、元のファイルを書き換える、指定したパスに書き出す、書き換え前に出力結果を確認してから書き出すこともできます。

## インストール

`ai-text-shaper`をインストールするには、以下のリンクからバイナリをダウンロードしてください：

[https://github.com/ytka/ai-text-shaper/releases](https://github.com/ytka/ai-text-shaper/releases)

## OpenAI APIのAPIキーを設定
このツールはOpenAI APIを利用するため APIキーが必要です。
ホームディレクトリに `.ai-text-shaper-apikey` という名前のファイルを作成し、OpenAI APIキーを書き込んでください。

## 使い方

`ai-text-shaper`の一般的な使用パターンは以下の通りです：

```sh
ai-text-shaper [オプション] [入力ファイル...]
```

入力ファイルは一つもしくは複数のファイルを指定できます。入力ファイルが指定されない場合、標準入力から読み取ります。

### オプション

#### プロンプトオプション

- `-p, --prompt string`
   - プロンプトのテキストを指定します。

- `-P, --prompt-path string`
   - プロンプトファイル（テキストファイル）のパスを指定します。このファイルから読み取った文字列をプロンプトとして使用します。
- '-m, --model string'
   - 使用するChat用モデルを指定します。デフォルトは `gpt-4o` です。

#### 出力オプション

- `-v, --verbose`
   - 詳細出力を有効にします。

- `-s, --silent`
   - すべての出力を抑制します。

- `-d, --diff`
   - 通常出力に加えて、入力と出力のテキストの差分を表示します。

#### ファイル書き込みオプション

- `-r, --rewrite`
   - 結果で入力ファイルを書き換えます。

- `-o, --outpath string`
   - 出力ファイルのパスを指定します。

- `-f, --use-first-code-block`
   - 出力テキストにコードブロックが含まれる場合、最初のコードブロックを出力として使用します。

- `-c, --confirm`
   - ファイルに書き込む前に書き込んでよいか確認を求めます。

## 使用例

### 基本的な使用方法

プロンプトをコマンドラインから与えるには：
```sh
ai-text-shaper -p "プロンプトのテキスト" /path/to/inputfile.txt
```

### プロンプトファイルを使用

プロンプトをファイルから与えるには：
```sh
ai-text-shaper -P /path/to/promptfile.txt /path/to/inputfile.txt
```

### 詳細出力

詳細出力を有効にするには：

```sh
ai-text-shaper -v /path/to/inputfile.txt
```

### 出力の抑制

全ての出力を抑制するには：

```sh
ai-text-shaper -s /path/to/inputfile.txt
```

### 差分の表示

入力と出力のテキストの差分を表示するには：

```sh
ai-text-shaper -d /path/to/inputfile.txt
```

### ファイルへの書き込み

結果を特定の出力ファイルに書き込むには：

```sh
ai-text-shaper -o /path/to/outputfile.txt /path/to/inputfile.txt
```

### 入力ファイルの書き換え

入力ファイルを結果で書き換える(rewrite)には：

```sh
ai-text-shaper -r /path/to/inputfile.txt
```

### 最初のコードブロックの使用

出力テキストの最初のコードブロックを使用するには：

```sh
ai-text-shaper -f /path/to/inputfile.txt
```

### 書き込み前の確認

ファイルに書き込む前に確認するには：

```sh
ai-text-shaper -c /path/to/inputfile.txt
```

## ライセンス

このプロジェクトは[MITライセンス](link_to_license)の下でライセンスされています。
