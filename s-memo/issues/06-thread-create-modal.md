# 06-thread-create-modal: スレッド作成モーダル

## 概要
新しいスレッドを作成するためのモーダルダイアログを実装する

## 作業内容
1. モーダルコンポーネントの作成
2. スレッド作成フォーム
3. フォームバリデーション
4. API連携
5. 作成後のリダイレクト処理
6. モーダル関連コンポーネントの単体テスト実装

## 成果物
- components/Modal.tsx (汎用モーダル)
- components/CreateThreadModal.tsx (スレッド作成モーダル)
- components/ThreadForm.tsx (スレッドフォーム)
- api/threads.ts (API処理)
- __tests__/components/Modal.test.tsx (モーダルテスト)
- __tests__/components/CreateThreadModal.test.tsx (スレッド作成モーダルテスト)
- __tests__/components/ThreadForm.test.tsx (スレッドフォームテスト)

## 想定作業量
- コード差分: 約280行
- 作業時間: 3-4時間

## 依存関係
- 01-project-setup
- 02-database-schema
- 05-thread-list-page

## 完了条件
- モーダルが正常に表示される
- スレッドが正常に作成される
- バリデーションが適切に動作する
- 作成後にモーダルが閉じる
- 作成したスレッドが一覧に反映される