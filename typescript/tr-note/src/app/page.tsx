import Link from "next/link";
import styles from "./page.module.css";
import { getAllTrainingRecords } from "../actions/training";

export default async function Home() {
  const trainingRecords = await getAllTrainingRecords();

  return (
    <div className={styles.page}>
      <main className={styles.main}>
        <h1 className={styles.title}>トレーニング記録</h1>

        <div className={styles.addButtonContainer}>
          <Link href="/add" className={styles.addButton}>
            新規記録を追加
          </Link>
        </div>

        {trainingRecords.length === 0 ? (
          <p className={styles.emptyMessage}>
            トレーニング記録がありません。「新規記録を追加」ボタンから記録を追加してください。
          </p>
        ) : (
          <div className={styles.recordList}>
            {trainingRecords.map((record) => (
              <div key={record.id} className={styles.recordCard}>
                <div className={styles.recordHeader}>
                  <h2 className={styles.recordTitle}>{record.title}</h2>
                  <span className={styles.recordDate}>{record.date}</span>
                </div>
                <p className={styles.recordDescription}>
                  {record.description.length > 100
                    ? `${record.description.substring(0, 100)}...`
                    : record.description}
                </p>
                <div className={styles.recordMeta}>
                  <span>種目数: {record.exercises.length}</span>
                </div>
                <div className={styles.recordActions}>
                  <Link
                    href={`/record/${record.id}`}
                    className={styles.viewButton}
                  >
                    詳細を見る
                  </Link>
                  <Link
                    href={`/delete/${record.id}`}
                    className={styles.deleteButton}
                  >
                    削除
                  </Link>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>
    </div>
  );
}
