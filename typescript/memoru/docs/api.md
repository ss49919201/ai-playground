# API 設計

## 概要

Memoru アプリケーションのバックエンド API は、Cloudflare Workers を使用して実装されます。API は RESTful な設計に従い、JSON フォーマットでデータをやり取りします。

## ベース URL

```
https://api.memoru.app
```

開発環境では以下の URL を使用します：

```
http://localhost:8787
```

## API エンドポイント

### 1. メモ一覧取得

#### リクエスト

```
GET /api/memos
```

#### クエリパラメータ

| パラメータ | 型     | 必須 | 説明                                                                                                        |
| ---------- | ------ | ---- | ----------------------------------------------------------------------------------------------------------- |
| page       | number | 任意 | ページ番号（デフォルト: 1）                                                                                 |
| limit      | number | 任意 | 1 ページあたりの件数（デフォルト: 10）                                                                      |
| sort       | string | 任意 | ソート順（created_at:desc, created_at:asc, updated_at:desc, updated_at:asc）（デフォルト: created_at:desc） |

#### レスポンス

**成功時 (200 OK)**

```json
{
  "memos": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "メモのタイトル",
      "content": "メモの内容",
      "createdAt": "2025-04-01T10:00:00.000Z",
      "updatedAt": "2025-04-01T10:00:00.000Z"
    },
    ...
  ],
  "pagination": {
    "total": 25,
    "page": 1,
    "limit": 10,
    "totalPages": 3
  }
}
```

**エラー時 (4xx/5xx)**

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "エラーメッセージ"
  }
}
```

### 2. メモ詳細取得

#### リクエスト

```
GET /api/memos/:id
```

#### パスパラメータ

| パラメータ | 型     | 必須 | 説明      |
| ---------- | ------ | ---- | --------- |
| id         | string | 必須 | メモの ID |

#### レスポンス

**成功時 (200 OK)**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "メモのタイトル",
  "content": "メモの内容",
  "createdAt": "2025-04-01T10:00:00.000Z",
  "updatedAt": "2025-04-01T10:00:00.000Z"
}
```

**エラー時 (4xx/5xx)**

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "エラーメッセージ"
  }
}
```

### 3. メモ作成

#### リクエスト

```
POST /api/memos
```

#### リクエストボディ

```json
{
  "title": "新しいメモのタイトル",
  "content": "新しいメモの内容"
}
```

| フィールド | 型     | 必須 | 説明                            |
| ---------- | ------ | ---- | ------------------------------- |
| title      | string | 必須 | メモのタイトル（最大 100 文字） |
| content    | string | 必須 | メモの内容（最大 10000 文字）   |

#### レスポンス

**成功時 (201 Created)**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "新しいメモのタイトル",
  "content": "新しいメモの内容",
  "createdAt": "2025-04-01T10:00:00.000Z",
  "updatedAt": "2025-04-01T10:00:00.000Z"
}
```

**エラー時 (4xx/5xx)**

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "エラーメッセージ",
    "details": {
      "title": [
        "タイトルは必須です",
        "タイトルは100文字以内で入力してください"
      ],
      "content": ["内容は必須です"]
    }
  }
}
```

### 4. メモ更新

#### リクエスト

```
PUT /api/memos/:id
```

#### パスパラメータ

| パラメータ | 型     | 必須 | 説明      |
| ---------- | ------ | ---- | --------- |
| id         | string | 必須 | メモの ID |

#### リクエストボディ

```json
{
  "title": "更新されたメモのタイトル",
  "content": "更新されたメモの内容"
}
```

| フィールド | 型     | 必須 | 説明                            |
| ---------- | ------ | ---- | ------------------------------- |
| title      | string | 必須 | メモのタイトル（最大 100 文字） |
| content    | string | 必須 | メモの内容（最大 10000 文字）   |

#### レスポンス

**成功時 (200 OK)**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "更新されたメモのタイトル",
  "content": "更新されたメモの内容",
  "createdAt": "2025-04-01T10:00:00.000Z",
  "updatedAt": "2025-04-01T10:30:00.000Z"
}
```

**エラー時 (4xx/5xx)**

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "エラーメッセージ",
    "details": {
      "title": [
        "タイトルは必須です",
        "タイトルは100文字以内で入力してください"
      ],
      "content": ["内容は必須です"]
    }
  }
}
```

### 5. メモ削除

#### リクエスト

```
DELETE /api/memos/:id
```

#### パスパラメータ

| パラメータ | 型     | 必須 | 説明      |
| ---------- | ------ | ---- | --------- |
| id         | string | 必須 | メモの ID |

#### レスポンス

**成功時 (204 No Content)**

レスポンスボディなし

**エラー時 (4xx/5xx)**

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "エラーメッセージ"
  }
}
```

## エラーコード

| コード                | HTTP ステータス | 説明                           |
| --------------------- | --------------- | ------------------------------ |
| INVALID_REQUEST       | 400             | リクエストの形式が不正         |
| VALIDATION_ERROR      | 400             | バリデーションエラー           |
| NOT_FOUND             | 404             | リソースが見つからない         |
| METHOD_NOT_ALLOWED    | 405             | 許可されていない HTTP メソッド |
| INTERNAL_SERVER_ERROR | 500             | サーバー内部エラー             |

## データモデル

### Memo モデル

```typescript
interface Memo {
  id: string; // UUID v4
  title: string; // メモのタイトル（最大100文字）
  content: string; // メモの内容（最大10000文字）
  createdAt: string; // ISO 8601形式の日時文字列
  updatedAt: string; // ISO 8601形式の日時文字列
}
```

## データベーススキーマ

Cloudflare D1（SQLite 互換）のスキーマ定義：

```sql
CREATE TABLE memos (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

-- インデックス
CREATE INDEX idx_memos_created_at ON memos(created_at);
CREATE INDEX idx_memos_updated_at ON memos(updated_at);
```

## セキュリティ

- すべての API リクエストは HTTPS 経由で行われます
- リクエストレート制限を実装して、DoS 攻撃を防止します
- 入力データは厳格にバリデーションされます
- エラーメッセージは適切に抽象化され、システム内部の詳細を漏らしません

## パフォーマンス最適化

- Cloudflare Workers のエッジでの実行により、低レイテンシーを実現
- レスポンスにはキャッシュヘッダーが適切に設定されます
- 大量のデータ取得にはページネーションを使用
- データベースクエリにはインデックスを適切に設定

## API 拡張計画（将来的な機能）

- 認証機能の追加（JWT 認証）
- メモの検索機能
- タグ付け機能
- メモの共有機能
- 画像アップロード機能
