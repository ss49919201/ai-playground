"use client";

import { useRef, useState } from "react";
import styles from "../page.module.css";

interface AddFormProps {
  onSubmit: (formData: FormData) => Promise<{ error?: string } | undefined>;
}

export default function AddForm({ onSubmit }: AddFormProps) {
  const formRef = useRef<HTMLFormElement>(null);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);

    try {
      const formData = new FormData(e.currentTarget);
      const result = await onSubmit(formData);

      if (result?.error) {
        setError(result.error);
        setIsSubmitting(false);
      }
    } catch {
      setError("エラーが発生しました。もう一度お試しください。");
      setIsSubmitting(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit} ref={formRef}>
      {error && <div className={styles.errorMessage}>{error}</div>}

      <div className={styles.formGroup}>
        <label htmlFor="title" className={styles.label}>
          タイトル *
        </label>
        <input
          id="title"
          name="title"
          type="text"
          className={styles.input}
          placeholder="例: 胸トレーニング、脚の日など"
          required
          disabled={isSubmitting}
        />
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="date" className={styles.label}>
          日付 *
        </label>
        <input
          id="date"
          name="date"
          type="date"
          className={styles.input}
          defaultValue={new Date().toISOString().split("T")[0]}
          required
          disabled={isSubmitting}
        />
      </div>

      <div className={styles.formGroup}>
        <label htmlFor="description" className={styles.label}>
          メモ
        </label>
        <textarea
          id="description"
          name="description"
          className={styles.textarea}
          placeholder="トレーニングに関するメモを入力してください"
          disabled={isSubmitting}
        />
      </div>

      <button
        type="submit"
        className={styles.submitButton}
        disabled={isSubmitting}
      >
        {isSubmitting ? "保存中..." : "保存"}
      </button>
    </form>
  );
}
