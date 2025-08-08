"use client";

import { updatePost } from "../actions";
import { useState } from "react";
import styles from "./PostForm.module.css";

type PostFormProps = {
  postId: string;
  initialTitle: string;
  initialBody: string;
  onCancel: () => void;
  onSave: () => void;
};

export default function PostForm({ postId, initialTitle, initialBody, onCancel, onSave }: PostFormProps) {
  const [title, setTitle] = useState(initialTitle);
  const [body, setBody] = useState(initialBody);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    
    try {
      await updatePost(postId, title, body);
      onSave();
    } catch (error) {
      console.error("Failed to update post:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      <div className={styles.field}>
        <label htmlFor="title" className={styles.label}>
          タイトル
        </label>
        <input
          id="title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className={styles.input}
          required
        />
      </div>
      
      <div className={styles.field}>
        <label htmlFor="body" className={styles.label}>
          内容
        </label>
        <textarea
          id="body"
          value={body}
          onChange={(e) => setBody(e.target.value)}
          className={styles.textarea}
          rows={4}
          required
        />
      </div>
      
      <div className={styles.buttons}>
        <button
          type="button"
          onClick={onCancel}
          className={styles.cancelButton}
          disabled={isSubmitting}
        >
          キャンセル
        </button>
        <button
          type="submit"
          className={styles.saveButton}
          disabled={isSubmitting}
        >
          {isSubmitting ? "保存中..." : "保存"}
        </button>
      </div>
    </form>
  );
}