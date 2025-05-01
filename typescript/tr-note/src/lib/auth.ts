"use client";

import { useEffect, useState } from "react";
import { Session } from "@auth/core";

// セッション取得用フック
export function useSession(): { session: Session | null; loading: boolean } {
  const [session, setSession] = useState<Session | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchSession() {
      try {
        // セッション情報を取得
        const res = await fetch("/api/auth/session");
        if (res.ok) {
          const data = await res.json();
          setSession(data.session || null);
        } else {
          setSession(null);
        }
      } catch (error) {
        console.error("セッション取得エラー:", error);
        setSession(null);
      } finally {
        setLoading(false);
      }
    }

    fetchSession();
  }, []);

  return { session, loading };
}

// サインアウト用関数
export async function signOut(): Promise<void> {
  await fetch("/api/auth/signout", { method: "POST" });
  window.location.href = "/";
}