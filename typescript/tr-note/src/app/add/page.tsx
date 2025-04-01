"use client";

import { useState } from "react";
import { useTraining } from "../../contexts/TrainingContext";
import { useRouter } from "next/navigation";
import Link from "next/link";
import styles from "../page.module.css";

export default function AddTrainingRecord() {
  const router = useRouter();
  const { addTrainingRecord } = useTraining();
  const [title, setTitle] = useState("");
  const [date, setDate] = useState(new Date().toISOString().split("T")[0]);
  const [description, setDescription] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!title || !date) {
      alert("タイトルと日付は必須です");
      return;
    }

    addTrainingRecord({
      title,
      date,
      description,
      exercises: [],
    });

    console.log("トレーニング記録が追加されました:", {
      title,
      date,
      description,
    });
    router.push("/");
  };

  return (
    <div className={styles.page}>
      <main className={styles.main}>
        <Link href="/" className={styles.backLink}>
          ← 戻る
        </Link>

        <h1 className={styles.title}>トレーニング記録の追加</h1>

        <form className={styles.form} onSubmit={handleSubmit}>
          <div className={styles.formGroup}>
            <label htmlFor="title" className={styles.label}>
              タイトル *
            </label>
            <input
              id="title"
              type="text"
              className={styles.input}
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="例: 胸トレーニング、脚の日など"
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label htmlFor="date" className={styles.label}>
              日付 *
            </label>
            <input
              id="date"
              type="date"
              className={styles.input}
              value={date}
              onChange={(e) => setDate(e.target.value)}
              required
            />
          </div>

          <div className={styles.formGroup}>
            <label htmlFor="description" className={styles.label}>
              メモ
            </label>
            <textarea
              id="description"
              className={styles.textarea}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="トレーニングに関するメモを入力してください"
            />
          </div>

          <button type="submit" className={styles.submitButton}>
            保存
          </button>
        </form>
      </main>
    </div>
  );
}
