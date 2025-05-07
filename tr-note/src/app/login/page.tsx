import Link from "next/link";
import styles from "../page.module.css";
import LoginForm from "./LoginForm";

export default function LoginPage() {
  return (
    <div className={styles.page}>
      <main className={styles.main}>
        <Link href="/" className={styles.backLink}>
          ← トップページに戻る
        </Link>
        
        <h1 className={styles.title}>ログイン</h1>
        
        <LoginForm />
        
        <div className={styles.formGroup} style={{ textAlign: "center", marginTop: "24px" }}>
          <p>アカウントをお持ちでない方は
            <Link href="/signup" style={{ textDecoration: "underline", marginLeft: "8px" }}>
              新規登録
            </Link>
          </p>
        </div>
      </main>
    </div>
  );
}
