"use client";

import { useTransition } from "react";
import { useRouter } from "next/navigation";
import styles from "../../page.module.css";

interface DeleteFormProps {
  onDelete: () => Promise<void>;
  recordId: string;
}

export default function DeleteForm({ onDelete, recordId }: DeleteFormProps) {
  const router = useRouter();
  const [isPending, startTransition] = useTransition();

  const handleDelete = () => {
    startTransition(async () => {
      await onDelete();
    });
  };

  const handleCancel = () => {
    router.push(`/record/${recordId}`);
  };

  return (
    <div className={styles.deleteActions}>
      {isPending ? (
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
  );
}
