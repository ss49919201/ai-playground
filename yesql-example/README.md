# YeSQL Account System

YeSQLの思想に従って実装された口座管理システムのバックエンドAPIです。セッションベースの認証機能を搭載しています。

## 機能

### 認証機能
- ユーザー登録
- ログイン/ログアウト
- セッション管理
- パスワードハッシュ化（bcrypt）

### 口座管理機能
- 口座作成
- 入金 (Deposit)
- 出金 (Withdraw) 
- 送金 (Transfer)
- 口座情報取得
- 取引履歴取得
- ユーザー別口座管理

## アーキテクチャ

このプロジェクトはYeSQL（Yet another SQL）の思想に従い、SQLクエリをGoコードから分離し、
`.sql`ファイルとして管理しています。

### ディレクトリ構造

```
yesql-account-system/
├── cmd/server/          # HTTPサーバーのエントリーポイント
├── internal/
│   ├── account/         # アカウント関連のビジネスロジック
│   ├── auth/           # 認証関連のビジネスロジック
│   ├── db/             # データベース接続管理
│   └── yesql/          # YeSQLクエリローダー
├── sql/
│   ├── schema.sql      # データベーススキーマ
│   └── queries/        # SQLクエリファイル
│       ├── accounts.sql
│       ├── transactions.sql
│       ├── users.sql
│       ├── sessions.sql
│       └── user_accounts.sql
└── bin/                # ビルド後の実行ファイル
```

## API エンドポイント

### 認証 (認証不要)
- `POST /auth/register` - ユーザー登録
- `POST /auth/login` - ログイン

### 保護されたAPI (認証必要 - `/api` プレフィックス)

#### 認証管理
- `POST /api/logout` - ログアウト

#### 口座管理
- `POST /api/accounts` - 口座作成
- `GET /api/accounts` - ユーザーの口座一覧
- `GET /api/accounts/{id}` - 口座詳細取得
- `GET /api/accounts/{id}/transactions` - 口座の取引履歴

#### 取引
- `POST /api/deposit` - 入金
- `POST /api/withdraw` - 出金  
- `POST /api/transfer` - 送金

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

### 1. ユーザー登録
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "tanaka", "email": "tanaka@example.com", "password": "password123"}'
```

### 2. ログイン
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "tanaka", "password": "password123"}' \
  -c cookies.txt
```

### 3. 口座作成（認証必要）
```bash
curl -X POST http://localhost:8080/api/accounts \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"account_id": "acc1", "account_name": "田中太郎の口座", "initial_deposit": 10000}'
```

### 4. 入金（認証必要）
```bash
curl -X POST http://localhost:8080/api/deposit \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"account_id": "acc1", "amount": 5000, "description": "給与振込"}'
```

### 5. 送金（認証必要）
```bash
curl -X POST http://localhost:8080/api/transfer \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"from_account_id": "acc1", "to_account_id": "acc2", "amount": 3000, "description": "家賃支払い"}'
```

### 6. ログアウト
```bash
curl -X POST http://localhost:8080/api/logout \
  -b cookies.txt
```

## 技術仕様

- **言語**: Go 1.21+
- **データベース**: SQLite3
- **HTTPルーター**: Gorilla Mux
- **認証**: セッションベース + bcryptパスワードハッシュ化
- **アーキテクチャ**: YeSQL (SQL分離)
- **テスト**: Go標準テストパッケージ

## セキュリティ機能

- **パスワードハッシュ化**: bcryptを使用した安全なパスワード保存
- **セッション管理**: 有効期限付きセッション（24時間）
- **認証ミドルウェア**: 保護されたエンドポイントへのアクセス制御
- **ユーザー別アクセス制御**: ユーザーは自分の口座のみアクセス可能
- **HTTPOnly Cookie**: XSS攻撃対策