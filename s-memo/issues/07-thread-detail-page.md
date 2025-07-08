# 07-thread-detail-page: スレッド詳細ページ

## 概要
選択したスレッドの詳細とメモ一覧を表示するページを作成する

## 作業内容
1. スレッド詳細ページの作成
2. スレッドヘッダーコンポーネント
3. メモ一覧表示
4. 新規メモ作成エリア
5. 時系列でのメモ表示
6. スレッド詳細関連コンポーネントの単体テスト実装

## 成果物
- pages/thread/[id].tsx (スレッド詳細ページ)
- components/ThreadHeader.tsx (スレッドヘッダー)
- components/MemoList.tsx (メモ一覧)
- components/MemoItem.tsx (メモアイテム)
- components/NewMemoForm.tsx (新規メモフォーム)
- __tests__/pages/thread/[id].test.tsx (スレッド詳細ページテスト)
- __tests__/components/MemoList.test.tsx (メモ一覧テスト)
- __tests__/components/MemoItem.test.tsx (メモアイテムテスト)

## 想定作業量
- コード差分: 約400行
- 作業時間: 5-6時間

## 依存関係
- 01-project-setup
- 02-database-schema
- 05-thread-list-page

## 完了条件
- スレッド詳細ページが表示される
- スレッド情報が正しく表示される
- メモ一覧が時系列で表示される
- 新規メモ作成フォームが機能する
- レスポンシブデザインが適用されている