"use client";

import { useTraining } from "../../../contexts/TrainingContext";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";
import styles from "../../page.module.css";
import { useEffect, useState } from "react";
import { TrainingRecord } from "../../../types";

export default function TrainingRecordDetail() {
  const params = useParams();
  const router = useRouter();
  const { getTrainingRecord } = useTraining();
  const [record, setRecord] = useState<TrainingRecord | undefined>(undefined);

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
          ← 戻る
        </Link>

        <div className={styles.recordDetail}>
          <div className={styles.detailHeader}>
            <h1 className={styles.detailTitle}>{record.title}</h1>
            <span className={styles.detailDate}>{record.date}</span>
          </div>

          {record.description && (
            <div className={styles.detailDescription}>{record.description}</div>
          )}

          <h2>トレーニング種目</h2>
          {record.exercises.length === 0 ? (
            <p>種目が登録されていません。</p>
          ) : (
            <div className={styles.exerciseList}>
              {record.exercises.map((exercise) => (
                <div key={exercise.id} className={styles.exerciseItem}>
                  <h3 className={styles.exerciseName}>{exercise.name}</h3>
                  {exercise.sets.length === 0 ? (
                    <p>セットが登録されていません。</p>
                  ) : (
                    <div className={styles.setList}>
                      {exercise.sets.map((set, index) => (
                        <div key={set.id} className={styles.setItem}>
                          <div className={styles.setNumber}>#{index + 1}</div>
                          <div className={styles.setDetail}>
                            <div>{set.weight} kg</div>
                            <div>{set.reps} 回</div>
                          </div>
                          {set.notes && (
                            <div className={styles.setNotes}>{set.notes}</div>
                          )}
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>

        <div className={styles.recordActions}>
          <Link href={`/delete/${record.id}`} className={styles.deleteButton}>
            このトレーニング記録を削除
          </Link>
        </div>
      </main>
    </div>
  );
}
