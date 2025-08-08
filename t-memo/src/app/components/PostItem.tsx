"use client";

import { useState } from "react";
import PostForm from "./PostForm";
import styles from "../page.module.css";

type Post = {
  id: string;
  title: string;
  body: string;
};

type PostItemProps = {
  post: Post;
};

export default function PostItem({ post }: PostItemProps) {
  const [isEditing, setIsEditing] = useState(false);

  const handleEdit = () => {
    setIsEditing(true);
  };

  const handleCancel = () => {
    setIsEditing(false);
  };

  const handleSave = () => {
    setIsEditing(false);
  };

  if (isEditing) {
    return (
      <div className={styles.postItem}>
        <PostForm
          postId={post.id}
          initialTitle={post.title}
          initialBody={post.body}
          onCancel={handleCancel}
          onSave={handleSave}
        />
      </div>
    );
  }

  return (
    <div className={styles.postItem}>
      <div className={styles.postHeader}>
        <h2 className={styles.postTitle}>{post.title}</h2>
        <button
          onClick={handleEdit}
          className={styles.editButton}
        >
          編集
        </button>
      </div>
      <p className={styles.postBody}>{post.body}</p>
    </div>
  );
}