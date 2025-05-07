"use client";

import { useState, FormEvent } from "react";
import { useRouter } from "next/navigation";
import styles from "../page.module.css";
import { createUser, getUserByEmail } from "../../db/user";
import { signIn } from "../../auth";

export default function SignupForm() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const validateEmail = (email: string) => {
    const re = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    return re.test(email);
  };

  const validatePassword = (password: string) => {
    return password.length >= 8;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError("");

    // バリデーション
    if (!name.trim()) {
      setError("名前を入力してください");
      return;
    }

    if (!email.trim()) {
      setError("メールアドレスを入力してください");
      return;
    }

    if (!validateEmail(email)) {
      setError("有効なメールアドレスを入力してください");
      return;
    }

    if (!password) {
      setError("パスワードを入力してください");
      return;
    }

    if (!validatePassword(password)) {
      setError("パスワードは8文字以上である必要があります");
      return;
    }

    if (password !== confirmPassword) {
      setError("パスワードが一致しません");
      return;
    }

    setIsLoading(true);

    try {
      // メールアドレスの重複チェック
      const existingUser = await getUserByEmail(email);
      if (existingUser) {
        setError("このメールアドレスは既に登録されています");
        setIsLoading(false);
        return;
      }

      // ユーザー作成
      await createUser(email, password, name);

      // 自動ログイン
      const signInResult = await signIn("credentials", {
        email,
        password,
        redirect: false,
      });

      if (signInResult?.error) {
        setError("ログインに失敗しました。もう一度お試しください。");
        setIsLoading(false);
        return;
      }

      // トップページにリダイレクト
      router.push("/");
      router.refresh();
    } catch (err) {
      console.error("登録エラー:", err);
      setError("アカウント登録中にエラーが発生しました。もう一度お試しください。");
      setIsLoading(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      {error && <div className={styles.errorMessage}>{error}</div>}

      <div className={styles.formGroup}>
        <label htmlFor="name" className={styles.label}>
          名前
        </label>
        <input
          id="name"
          type="text"
          className={styles.input}
          value={name}
          onChange={(e) => setName(e.target.value)}
          disabled={isLoading}
          required
        />
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="email" className={styles.label}>
          メールアドレス
        </label>
        <input
          id="email"
          type="email"
          className={styles.input}
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          disabled={isLoading}
          required
        />
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="password" className={styles.label}>
          パスワード
        </label>
        <input
          id="password"
          type="password"
          className={styles.input}
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          disabled={isLoading}
          required
          minLength={8}
        />
        <small style={{ color: "rgba(var(--gray-rgb), 0.7)" }}>
          8文字以上で入力してください
        </small>
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="confirmPassword" className={styles.label}>
          パスワード（確認）
        </label>
        <input
          id="confirmPassword"
          type="password"
          className={styles.input}
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          disabled={isLoading}
          required
        />
      </div>

      <button
        type="submit"
        className={styles.submitButton}
        disabled={isLoading}
      >
        {isLoading ? "登録中..." : "登録する"}
      </button>
    </form>
  );
}
