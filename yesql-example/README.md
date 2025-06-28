# YeSQL Account System

YeSQLの思想に従って実装された口座管理システムのバックエンドAPIです。

## 機能

- 口座作成
- 入金 (Deposit)
- 出金 (Withdraw) 
- 送金 (Transfer)
- 口座情報取得
- 取引履歴取得

## アーキテクチャ

このプロジェクトはYeSQL（Yet another SQL）の思想に従い、SQLクエリをGoコードから分離し、
`.sql`ファイルとして管理しています。

### ディレクトリ構造

```
yesql-account-system/
├── cmd/server/          # HTTPサーバーのエントリーポイント
├── internal/
│   ├── account/         # アカウント関連のビジネスロジック
│   ├── db/             # データベース接続管理
│   └── yesql/          # YeSQLクエリローダー
├── sql/
│   ├── schema.sql      # データベーススキーマ
│   └── queries/        # SQLクエリファイル
│       ├── accounts.sql
│       └── transactions.sql
└── bin/                # ビルド後の実行ファイル
```

## API エンドポイント

### 口座管理
- `POST /accounts` - 口座作成
- `GET /accounts` - 全口座一覧
- `GET /accounts/{id}` - 口座詳細取得
- `GET /accounts/{id}/transactions` - 口座の取引履歴

### 取引
- `POST /deposit` - 入金
- `POST /withdraw` - 出金  
- `POST /transfer` - 送金

### ヘルスチェック
- `GET /health` - サーバー状態確認

## 実行方法

1. プロジェクトのビルド:
```bash
go build -o bin/server cmd/server/*.go
```

2. サーバー起動:
```bash
./bin/server
```

サーバーは http://localhost:8080 で起動します。

## テスト実行

```bash
go test ./...
```

## 使用例

### 口座作成
```bash
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"account_id": "acc1", "account_name": "田中太郎", "initial_deposit": 10000}'
```

### 入金
```bash
curl -X POST http://localhost:8080/deposit \
  -H "Content-Type: application/json" \
  -d '{"account_id": "acc1", "amount": 5000, "description": "給与振込"}'
```

### 送金
```bash
curl -X POST http://localhost:8080/transfer \
  -H "Content-Type: application/json" \
  -d '{"from_account_id": "acc1", "to_account_id": "acc2", "amount": 3000, "description": "家賃支払い"}'
```

## 技術仕様

- **言語**: Go 1.21+
- **データベース**: SQLite3
- **HTTPルーター**: Gorilla Mux
- **アーキテクチャ**: YeSQL (SQL分離)
- **テスト**: Go標準テストパッケージ