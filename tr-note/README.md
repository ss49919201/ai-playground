# トレーニング記録アプリ

このアプリケーションは、トレーニング記録を管理するためのシンプルなWebアプリケーションです。トレーニングの日付、タイトル、内容、各種目のセット情報などを記録することができます。

## 技術スタック

- [Next.js](https://nextjs.org/)
- [TypeScript](https://www.typescriptlang.org/)
- [Cloudflare D1](https://developers.cloudflare.com/d1/) (SQLiteベースのデータベース)

## 機能

- トレーニング記録の追加
- トレーニング記録の閲覧
- トレーニング記録の削除
- トレーニング種目とセットの追加

## 開発環境のセットアップ

1. リポジトリをクローン

```bash
git clone <repository-url>
cd tr-note
```

2. 依存関係のインストール

```bash
npm install
```

3. ローカルD1データベースの作成

```bash
npm run d1:create
```

4. マイグレーションの実行

```bash
npm run d1:migrations:apply:local
```

5. 開発サーバーの起動

```bash
npm run dev
```

## ローカルD1の設定方法

トレーニング記録アプリをローカルで開発する場合、以下の手順でCloudflare D1を設定します。

### 1. D1データベースの作成

```bash
npm run d1:create
```

このコマンドはCloudflareアカウントに新しいD1データベースを作成します。wrangler.jsonc内の設定を更新するためのコマンドが表示されるので、それに従ってください。

### 2. マイグレーションの適用

次のコマンドでマイグレーションを適用します：

```bash
npm run d1:migrations:apply:local
```

これにより、`migrations`ディレクトリ内のSQLファイルが実行され、データベーススキーマが作成されます。

### 3. SQLの実行

特定のSQLコマンドを実行する場合：

```bash
npm run d1:execute:local --command="SELECT * FROM training_records"
```

もしくはSQLファイルを実行する場合：

```bash
npm run d1:execute:local --file=./queries/sample.sql
```

## Next.jsとD1の連携の仕組み

このアプリケーションでは、Next.jsのServer ActionsからD1データベースにアクセスするために次の方法を使用しています：

1. Cloudflare環境変数からD1インスタンスを取得：
```typescript
import { getCloudflareContext } from '@opennextjs/cloudflare';

export async function getDB() {
  const { env } = getCloudflareContext();
  return env.DB;
}
```

2. Server ActionsでD1を使用：
```typescript
"use server";
import { getDB } from "../db";

export async function getData() {
  const db = await getDB();
  const result = await db.prepare("SELECT * FROM my_table").all();
  return result.results;
}
```

## デプロイ

1. Cloudflareアカウントにログイン

```bash
npx wrangler login
```

2. D1データベースを作成（初回のみ）

```bash
npm run d1:create
```

3. マイグレーションを実行

```bash
npm run d1:migrations:apply
```

4. ビルドとデプロイ

```bash
npm run deploy
```

## データベース設計

トレーニング記録は以下の3つのテーブルで管理されています：

1. `training_records` - トレーニング記録の基本情報
2. `exercises` - トレーニング種目
3. `sets` - 各種目のセット情報

## Learn More

To learn more about the technologies used in this project:

- [Next.js Documentation](https://nextjs.org/docs)
- [Cloudflare D1 Documentation](https://developers.cloudflare.com/d1/)
- [Cloudflare Pages with Next.js](https://developers.cloudflare.com/pages/framework-guides/nextjs/)