import { authConfig, getAdapter } from "@/auth";
import NextAuth from "@auth/core";

// NextAuthハンドラーの作成
const handler = NextAuth({
  // 認証設定
  ...authConfig,
  // D1アダプターを設定
  adapter: getAdapter(),
});

// GETとPOSTリクエストを処理
export { handler as GET, handler as POST };