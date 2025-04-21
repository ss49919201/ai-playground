# データベース設計

## 概要
ユーザー認証に必要なデータベーステーブルの設計と既存テーブルの拡張

## タスク
- [ ] ユーザーテーブルの設計
  ```sql
  CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name TEXT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
  );
  ```
- [ ] セッションテーブルの設計
  ```sql
  CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    expires_at TEXT NOT NULL,
    created_at TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
  );
  ```
- [ ] 既存のトレーニング記録テーブルにユーザーIDカラム追加
  ```sql
  ALTER TABLE training_records ADD COLUMN user_id TEXT;
  ALTER TABLE training_records ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id);
  ```
- [ ] マイグレーションファイル作成 (`migrations/0001_auth.sql`)

## 担当者
未割り当て

## 依存関係
なし

## 見積時間
2時間