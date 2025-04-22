import { NextRequest, NextResponse } from "next/server";
import { auth } from "./auth";

// 認証が不要なパス
const publicPaths = ["/login", "/signup", "/"];

export async function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname;
  
  // 公開パスはそのまま通す
  if (publicPaths.some(path => pathname === path || pathname.startsWith(path + "/"))) {
    return NextResponse.next();
  }

  // セッション情報の取得
  const session = await auth();
  
  // 未認証の場合はログインページにリダイレクト
  if (!session) {
    const url = new URL("/login", request.url);
    url.searchParams.set("callbackUrl", encodeURI(request.url));
    return NextResponse.redirect(url);
  }
  
  return NextResponse.next();
}

// ミドルウェアを適用するパスを指定
export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico).*)"],
};