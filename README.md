# Perplexity API トークン平均計算ツール

このプログラムは、Perplexity AIのAPIを使用して複数の質問の応答に使われたトークン数を計算し、平均値を求めるツールです。

## 必要条件

- Go 1.16以上
- Perplexity AI APIキー

## 使い方

1. 環境変数に`SONARAPI_KEY`としてPerplexity APIキーを設定します:

```bash
export SONARAPI_KEY="your-api-key-here"
```

2. 質問を含むテキストファイルを作成します。質問はそれぞれ `---` で区切ります:

```
質問1テキスト（複数行可能）
---
質問2テキスト（複数行可能）
---
質問3テキスト（複数行可能）
```

3. プログラムを実行します:

```bash
# デフォルトのファイル名 (questions.txt) を使用する場合
go run main.go

# 異なるファイル名を指定する場合
go run main.go -file=カスタム質問ファイル.txt
```

## 機能

- テキストファイルから複数の質問を読み込み
- 複数の質問に対するPerplexity AIのレスポンスを取得
- 各質問の`completion_tokens`を表示
- すべての質問の`completion_tokens`の平均を計算して表示

## ファイル形式

質問ファイルは以下の形式で作成します:

- 各質問は `---` で区切ります
- 質問内の改行は保持されます
- 空白行や不要な空白は自動的に削除されます

## カスタマイズ

質問ファイルのフォーマットやパスは自由に変更できます。デフォルトでは、プログラムは同じディレクトリ内の `questions.txt` ファイルを読み込みます。
# inspect-sonar-api-average-completion-token
