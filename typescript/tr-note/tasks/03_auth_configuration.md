# 認証設定実装

## 概要
Auth.js設定ファイル作成とAPI Route実装

## タスク
- [ ] Auth.js（NextAuth）設定ファイル作成
  - [ ] `src/auth.ts` ファイルを作成
  - [ ] 認証プロバイダーの設定（Credentials Provider）
  - [ ] D1アダプターの設定
  - [ ] セッション設定
- [ ] 認証API Route作成
  - [ ] `src/app/api/auth/[...nextauth]/route.ts` ファイル作成
  - [ ] ハンドラー設定
- [ ] 環境変数設定
  - [ ] `AUTH_SECRET` の設定（本番環境用秘密鍵）
  - [ ] 必要に応じて `.env.local` ファイル作成（開発環境用）

## 担当者
未割り当て

## 依存関係
- タスク01: データベース設計
- タスク02: 認証パッケージ導入

## 見積時間
3時間