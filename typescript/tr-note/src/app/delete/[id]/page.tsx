"use client";

import { useTraining } from "../../../contexts/TrainingContext";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import styles from "../../page.module.css";
import { useEffect, useState } from "react";
import { TrainingRecord } from "../../../types";

export default function DeleteTrainingRecord() {
  const params = useParams();
  const router = useRouter();
  const { getTrainingRecord, deleteTrainingRecord } = useTraining();
  const [record, setRecord] = useState<TrainingRecord | undefined>(undefined);
  const [isDeleting, setIsDeleting] = useState(false);

  useEffect(() => {
    if (params.id) {
      const foundRecord = getTrainingRecord(params.id as string);
      setRecord(foundRecord);

      if (!foundRecord) {
        console.log(`ID: ${params.id} のトレーニング記録が見つかりません`);
        router.push("/");
      }
    }
  }, [params.id, getTrainingRecord, router]);

  const handleDelete = () => {
    if (record) {
      setIsDeleting(true);

      // 削除処理
      deleteTrainingRecord(record.id);
      console.log(`ID: ${record.id} のトレーニング記録が削除されました`);

      // 少し遅延させてから一覧ページに戻る
      setTimeout(() => {
        router.push("/");
      }, 1000);
    }
  };

  const handleCancel = () => {
    // 詳細ページに戻る
    if (record) {
      router.push(`/record/${record.id}`);
    } else {
      router.push("/");
    }
  };

  if (!record) {
    return (
      <div className={styles.page}>
        <main className={styles.main}>
          <p>読み込み中...</p>
        </main>
      </div>
    );
  }

  return (
    <div className={styles.page}>
      <main className={styles.main}>
        <Link href="/" className={styles.backLink}>
          ← ホームに戻る
        </Link>

        <div className={styles.deleteConfirmation}>
          <h1 className={styles.title}>トレーニング記録の削除</h1>

          <div className={styles.recordSummary}>
            <h2 className={styles.recordTitle}>{record.title}</h2>
            <span className={styles.recordDate}>{record.date}</span>
            {record.description && (
              <p className={styles.recordDescription}>
                {record.description.length > 100
                  ? `${record.description.substring(0, 100)}...`
                  : record.description}
              </p>
            )}
          </div>

          <div className={styles.warningMessage}>
            <p>このトレーニング記録を削除してもよろしいですか？</p>
            <p>この操作は取り消せません。</p>
          </div>

          <div className={styles.deleteActions}>
            {isDeleting ? (
              <p className={styles.deletingMessage}>削除中...</p>
            ) : (
              <>
                <button onClick={handleCancel} className={styles.cancelButton}>
                  キャンセル
                </button>
                <button onClick={handleDelete} className={styles.deleteButton}>
                  削除する
                </button>
              </>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
