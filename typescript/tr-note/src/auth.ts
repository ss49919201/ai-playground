import { D1Adapter } from "@auth/d1-adapter";
import { NextAuthConfig, Session, User, NextAuth } from "@auth/core";
import { getCloudflareContext } from "@opennextjs/cloudflare";
import Credentials from "@auth/core/providers/credentials";
import bcrypt from "bcrypt";
import { getUserByEmail, getUserPasswordHash } from "./db/user";

// 認証設定
export const authConfig: NextAuthConfig = {
  // セッション管理方法（JWT）
  session: {
    strategy: "jwt",
    maxAge: 30 * 24 * 60 * 60, // 30日
  },

  // 認証プロバイダー
  providers: [
    Credentials({
      id: "credentials",
      name: "Email & Password",
      credentials: {
        email: { label: "Email", type: "email" },
        password: { label: "Password", type: "password" }
      },
      // 認証ロジック
      async authorize(credentials) {
        if (!credentials?.email || !credentials?.password) {
          return null;
        }

        try {
          // メールアドレスでユーザーを検索
          const user = await getUserByEmail(credentials.email);
          if (!user) {
            return null;
          }

          // パスワードハッシュの取得と検証
          const passwordHash = await getUserPasswordHash(user.id);
          if (!passwordHash) {
            return null;
          }

          // パスワードの照合
          const isPasswordValid = await bcrypt.compare(credentials.password, passwordHash);
          if (!isPasswordValid) {
            return null;
          }

          // 認証成功時はユーザー情報を返す
          return {
            id: user.id,
            email: user.email,
            name: user.name,
          } as User;
        } catch (error) {
          console.error("認証エラー:", error);
          return null;
        }
      },
    }),
  ],

  // コールバック
  callbacks: {
    // ユーザー情報をセッションに含める
    async session({ session, token }: { session: Session; token: any }) {
      if (token.sub && session.user) {
        session.user.id = token.sub;
      }
      return session;
    },
  },

  // 認証関連ページのURLパス
  pages: {
    signIn: "/login",
    signOut: "/",
    error: "/login",
  },

  // クライアントサイドでも利用可能なフラグ
  useSecureCookies: process.env.NODE_ENV === "production",

  // デバッグモード
  debug: process.env.NODE_ENV === "development",
};

// D1アダプター設定用関数
export function getAdapter() {
  const { env } = getCloudflareContext();
  
  if (!env?.DB) {
    throw new Error("D1データベースが利用できません。wrangler.jsonc と open-next.config.ts で D1バインディングが正しく設定されているか確認してください。");
  }
  
  return D1Adapter(env.DB);
}

// Next.js App Routerでの使用のためにauth()関数をエクスポート
export const { auth, signIn, signOut } = NextAuth(authConfig);