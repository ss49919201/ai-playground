import Link from "next/link";
import styles from "../page.module.css";
import { redirect } from "next/navigation";
import { addTrainingRecord } from "../../actions/training";
import AddForm from "./AddForm";

export default function AddTrainingRecord() {
  async function handleAddRecord(formData: FormData) {
    "use server";

    const title = formData.get("title") as string;
    const date = formData.get("date") as string;
    const description = formData.get("description") as string;

    if (!title || !date) {
      return { error: "タイトルと日付は必須です" };
    }

    await addTrainingRecord({
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

    redirect("/");
  }

  return (
    <div className={styles.page}>
      <main className={styles.main}>
        <Link href="/" className={styles.backLink}>
          ← 戻る
        </Link>

        <h1 className={styles.title}>トレーニング記録の追加</h1>

        <AddForm onSubmit={handleAddRecord} />
      </main>
    </div>
  );
}
