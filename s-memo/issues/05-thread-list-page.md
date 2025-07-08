# 05-thread-list-page: スレッド一覧ページ

## 概要
スレッドの一覧を表示するメインページを作成する

## 作業内容
1. スレッド一覧ページの作成
2. スレッド一覧コンポーネント
3. 新規スレッド作成ボタン
4. スレッドカードコンポーネント
5. 検索・フィルタリング機能（基本）
6. スレッド一覧関連コンポーネントの単体テスト実装

## 成果物
- pages/index.tsx (メインページ)
- components/ThreadList.tsx (スレッド一覧)
- components/ThreadCard.tsx (スレッドカード)
- components/CreateThreadButton.tsx (作成ボタン)
- components/SearchBar.tsx (検索バー)
- __tests__/pages/index.test.tsx (メインページテスト)
- __tests__/components/ThreadList.test.tsx (スレッド一覧テスト)
- __tests__/components/ThreadCard.test.tsx (スレッドカードテスト)

## 想定作業量
- コード差分: 約350行
- 作業時間: 4-5時間

## 依存関係
- 01-project-setup
- 02-database-schema
- 03-auth-system

## 完了条件
- スレッド一覧が表示される
- スレッドカードが適切にレンダリングされる
- 新規作成ボタンが機能する
- 基本的な検索機能が動作する
- レスポンシブデザインが適用されている