# 08-memo-crud: メモのCRUD操作

## 概要
メモの作成、読み取り、更新、削除機能を実装する

## 作業内容
1. メモ作成機能
2. メモ編集機能
3. メモ削除機能
4. メモ表示機能
5. API エンドポイントの実装
6. メモCRUD関連コンポーネントの単体テスト実装

## 成果物
- api/memos/create.ts (メモ作成API)
- api/memos/update.ts (メモ更新API)
- api/memos/delete.ts (メモ削除API)
- components/MemoEditor.tsx (メモエディター)
- components/MemoActions.tsx (メモアクション)
- __tests__/api/memos/create.test.ts (メモ作成APIテスト)
- __tests__/api/memos/update.test.ts (メモ更新APIテスト)
- __tests__/api/memos/delete.test.ts (メモ削除APIテスト)
- __tests__/components/MemoEditor.test.tsx (メモエディターテスト)

## 想定作業量
- コード差分: 約450行
- 作業時間: 6-7時間

## 依存関係
- 01-project-setup
- 02-database-schema
- 07-thread-detail-page

## 完了条件
- メモが正常に作成される
- メモが正常に編集される
- メモが正常に削除される
- すべてのCRUD操作が正常に動作する
- エラーハンドリングが適切に実装されている