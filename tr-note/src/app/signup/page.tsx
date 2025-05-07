import Link from "next/link";
import styles from "../page.module.css";
import SignupForm from "./SignupForm";

export default function SignupPage() {
  return (
    <div className={styles.page}>
      <main className={styles.main}>
        <Link href="/" className={styles.backLink}>
          ← トップページに戻る
        </Link>
        
        <h1 className={styles.title}>新規ユーザー登録</h1>
        
        <SignupForm />
        
        <div className={styles.formGroup} style={{ textAlign: "center", marginTop: "24px" }}>
          <p>すでにアカウントをお持ちの方は
            <Link href="/login" style={{ textDecoration: "underline", marginLeft: "8px" }}>
              ログイン
            </Link>
          </p>
        </div>
      </main>
    </div>
  );
}
