# 03-auth-system: 認証システム基盤

## 概要
JWTを使用した認証システムの基盤を構築する

## 作業内容
1. JWT認証の実装
2. ログイン・ログアウト機能
3. 認証コンテキストの作成
4. 認証済みユーザーの情報管理
5. 保護されたルートの実装
6. 認証関連の単体テスト実装

## 成果物
- lib/auth.ts (JWT処理)
- contexts/AuthContext.tsx (認証コンテキスト)
- components/AuthGuard.tsx (認証ガード)
- hooks/useAuth.ts (認証フック)
- __tests__/lib/auth.test.ts (認証テスト)
- __tests__/hooks/useAuth.test.ts (認証フックテスト)

## 想定作業量
- コード差分: 約300行
- 作業時間: 4-5時間

## 依存関係
- 01-project-setup
- 02-database-schema

## 完了条件
- JWT認証が正常に動作する
- ログイン・ログアウトが機能する
- 認証状態が適切に管理される
- 未認証ユーザーは保護されたページにアクセスできない