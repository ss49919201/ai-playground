# 09-thread-management: スレッド管理機能

## 概要
スレッドの編集、削除、検索などの管理機能を実装する

## 作業内容
1. スレッド編集機能
2. スレッド削除機能
3. スレッド検索機能
4. スレッド並び替え機能
5. 関連するAPI実装
6. スレッド管理関連コンポーネントの単体テスト実装

## 成果物
- api/threads/update.ts (スレッド更新API)
- api/threads/delete.ts (スレッド削除API)
- components/ThreadEditModal.tsx (スレッド編集モーダル)
- components/ThreadDeleteModal.tsx (削除確認モーダル)
- components/ThreadSortControls.tsx (並び替え)
- __tests__/api/threads/update.test.ts (スレッド更新APIテスト)
- __tests__/api/threads/delete.test.ts (スレッド削除APIテスト)
- __tests__/components/ThreadEditModal.test.tsx (スレッド編集モーダルテスト)

## 想定作業量
- コード差分: 約350行
- 作業時間: 4-5時間

## 依存関係
- 01-project-setup
- 02-database-schema
- 05-thread-list-page
- 06-thread-create-modal

## 完了条件
- スレッドが正常に編集される
- スレッドが正常に削除される
- 検索機能が動作する
- 並び替えが機能する
- 削除確認ダイアログが表示される