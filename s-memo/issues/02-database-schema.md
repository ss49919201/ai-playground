# 02-database-schema: データベーススキーマ設計

## 概要
Cloudflare D1用のデータベーススキーマを設計・作成する

## 作業内容
1. データベーススキーマの設計
2. SQLマイグレーションファイルの作成
3. TypeScript型定義の作成
4. データベース接続設定

## 成果物
- schema.sql (初期マイグレーション)
- types/database.ts (型定義)
- lib/database.ts (DB接続設定)

## データベーステーブル
- users (ユーザー情報)
- threads (スレッド情報)
- memos (メモ情報)

## 想定作業量
- コード差分: 約150行
- 作業時間: 2-3時間

## 依存関係
- 01-project-setup

## 完了条件
- データベーススキーマが定義されている
- TypeScript型定義が完成している
- データベース接続設定が完了している