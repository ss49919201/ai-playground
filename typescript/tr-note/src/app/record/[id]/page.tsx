import Link from "next/link";
import styles from "../../page.module.css";
import { getTrainingRecord } from "../../../actions/training";
import { notFound } from "next/navigation";

export default async function TrainingRecordDetail({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  // Properly await params before accessing properties
  const { id } = await params;
  const record = await getTrainingRecord(id);

  if (!record) {
    console.log(`ID: ${id} のトレーニング記録が見つかりません`);
    notFound();
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
