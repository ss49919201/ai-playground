# Simple RAG System

Go言語で実装されたシンプルなRAG（Retrieval-Augmented Generation）システムです。ドキュメントを埋め込みベクトルに変換して保存し、質問に対して関連する文書を検索してLLMで回答を生成します。

## 特徴

- **モジュラー設計**: 設定、ドキュメント処理、ベクトル検索、LLM統合が分離された構造
- **JSON形式でのデータ保存**: シンプルなファイルベースのベクトルデータベース
- **Ollama LLM統合**: ローカルで動作するLLMサービスとの連携
- **外部埋め込みサービス対応**: sentence-transformersベースの埋め込み生成
- **CLIインターフェース**: コマンドラインから簡単に操作可能
- **包括的テストカバレッジ**: 全パッケージの単体テスト完備

## 必要な環境

### 依存サービス

1. **埋め込みサービス** (ポート 8000)
   - sentence-transformersを使用した埋め込み生成API
   - 例: `transformers`ライブラリを使ったHTTPサーバー

2. **Ollama LLMサービス** (ポート 11434)
   - ローカルLLMサービス
   - インストール: [Ollama公式サイト](https://ollama.ai/)

### Go環境
- Go 1.21以上

## インストール

```bash
# リポジトリをクローン
git clone <repository-url>
cd simple-rag

# 依存関係をインストール
go mod download

# ビルド
go build -o rag ./cmd/rag
```

## 設定

設定ファイル `config.yaml` をプロジェクトルートに配置します：

```yaml
server:
  port: 8080
  host: "localhost"

embedding:
  url: "http://localhost:8000"
  model: "all-MiniLM-L6-v2"
  batch_size: 32

llm:
  url: "http://localhost:11434"
  model: "llama2"
  temperature: 0.7
  max_tokens: 512

document:
  chunk_size: 512
  chunk_overlap: 50
  supported_formats: ["txt", "md", "text"]

vector_db:
  storage_path: "./data/vectors"
  similarity_threshold: 0.7

logging:
  level: "info"
  output: "stdout"
```

## 使用方法

### 1. ドキュメントの追加

テキストファイルをRAGシステムに追加します：

```bash
# 単一ファイルを追加
./rag add /path/to/document.txt

# ディレクトリ内の全ファイルを追加
./rag add /path/to/documents/
```

**対応ファイル形式:**
- `.txt` - プレーンテキスト
- `.md` - Markdown
- `.text` - テキストファイル

### 2. 質問と回答

追加したドキュメントに基づいて質問に回答します：

```bash
# 対話モード
./rag query

# 直接質問
./rag query "Goの特徴は何ですか？"
```

### 3. ドキュメント一覧表示

保存されているドキュメントを確認します：

```bash
./rag list
```

### 4. システムヘルスチェック

依存サービスの動作確認：

```bash
./rag health
```

## 使用例

### ドキュメント追加の例

```bash
# サンプルドキュメントを追加
./rag add ./data/documents/sample1.txt
./rag add ./data/documents/sample2.txt

# 追加結果確認
./rag list
```

### 質問例

```bash
# Go言語について質問
./rag query "Goプログラミング言語の特徴を教えて"

# 機械学習について質問  
./rag query "機械学習とは何ですか？"

# 対話モードで複数質問
./rag query
> Goの並行処理について詳しく教えて
> 機械学習のアルゴリズムにはどんなものがありますか？
> exit
```

## ディレクトリ構成

```
simple-rag/
├── cmd/rag/              # CLIアプリケーション
│   ├── main.go
│   ├── rag_system.go
│   └── utils.go
├── internal/             # 内部パッケージ
│   ├── config/           # 設定管理
│   ├── document/         # ドキュメント処理
│   ├── llm/             # LLMクライアント
│   └── vector/          # ベクトルDB・埋め込み
├── pkg/types/           # 共通データ型
├── data/                # データディレクトリ
│   ├── documents/       # 入力ドキュメント
│   └── vectors/         # ベクトルDB保存先
├── config.yaml          # 設定ファイル
└── README.md
```

## API仕様

### 埋め込みサービス API

RAGシステムが期待する埋め込みサービスのAPI仕様：

```bash
# ヘルスチェック
GET /health

# 埋め込み生成
POST /embeddings
Content-Type: application/json

{
  "texts": ["テキスト1", "テキスト2"],
  "model": "all-MiniLM-L6-v2"
}

# レスポンス
{
  "embeddings": [[0.1, 0.2, ...], [0.3, 0.4, ...]],
  "model": "all-MiniLM-L6-v2"
}
```

### Ollama API

標準のOllama API仕様に準拠：

```bash
# モデル一覧
GET /api/tags

# テキスト生成
POST /api/generate
{
  "model": "llama2",
  "prompt": "質問文",
  "stream": false,
  "options": {
    "temperature": 0.7,
    "num_predict": 512
  }
}
```

## トラブルシューティング

### 一般的な問題

**1. 埋め込みサービスに接続できない**
```bash
# ヘルスチェックで確認
./rag health

# サービスが起動しているか確認
curl http://localhost:8000/health
```

**2. Ollamaに接続できない**
```bash
# Ollamaの状態確認
ollama list

# モデルがインストールされているか確認
ollama pull llama2
```

**3. ドキュメントが追加されない**
- ファイル形式が対応しているか確認（txt, md, text）
- ファイルの読み取り権限があるか確認
- ファイルサイズが適切か確認

**4. 検索結果が期待通りでない**
- `similarity_threshold`の値を調整（config.yaml）
- `chunk_size`や`chunk_overlap`を調整
- より多くのドキュメントを追加

### ログ出力

詳細なログを確認するには設定ファイルでログレベルを変更：

```yaml
logging:
  level: "debug"  # info, debug, warn, error
  output: "stdout"
```

## 開発

### テスト実行

```bash
# 全テスト実行
go test ./...

# 特定パッケージのテスト
go test ./internal/config -v

# カバレッジ付きテスト
go test -cover ./...
```

### 新機能追加

1. 適切なパッケージに機能を追加
2. 対応するテストファイルを作成
3. テストが全てパスすることを確認
4. ドキュメントを更新

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。

## 貢献

プルリクエストやイシューの報告は歓迎します。大きな変更を行う前に、まずイシューを開いて相談してください。

## サポート

質問や問題がある場合は、GitHubのIssuesページでお知らせください。