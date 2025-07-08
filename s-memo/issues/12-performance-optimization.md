# 12-performance-optimization: パフォーマンス最適化

## 概要
アプリケーションのパフォーマンスを最適化する

## 作業内容
1. 画像最適化
2. コード分割の実装
3. キャッシュ戦略の実装
4. バンドルサイズの最適化
5. 読み込み時間の短縮

## 成果物
- next.config.js (Next.js最適化設定)
- components/LazyImage.tsx (遅延読み込み画像)
- hooks/useCache.ts (キャッシュフック)
- lib/performance.ts (パフォーマンス計測)

## 想定作業量
- コード差分: 約250行
- 作業時間: 3-4時間

## 依存関係
- 01-project-setup
- 05-thread-list-page
- 07-thread-detail-page

## 完了条件
- 初期読み込みが3秒以内
- 画像が適切に最適化される
- コード分割が正常に動作する
- キャッシュが適切に機能する
- パフォーマンス指標が改善される