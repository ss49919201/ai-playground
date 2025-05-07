import Link from "next/link";
import styles from "../../page.module.css";
import {
  getTrainingRecord,
  deleteTrainingRecord,
} from "../../../actions/training";
import { notFound, redirect } from "next/navigation";
import DeleteForm from "./DeleteForm";

export default async function DeleteTrainingRecord({
  params,
}: {
  params: { id: string };
}) {
  const record = await getTrainingRecord(params.id);

  if (!record) {
    console.log(`ID: ${params.id} のトレーニング記録が見つかりません`);
    notFound();
  }

  async function handleDelete() {
    "use server";
    await deleteTrainingRecord(params.id);
    console.log(`ID: ${params.id} のトレーニング記録が削除されました`);
    redirect("/");
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

          <DeleteForm onDelete={handleDelete} recordId={record.id} />
        </div>
      </main>
    </div>
  );
}
