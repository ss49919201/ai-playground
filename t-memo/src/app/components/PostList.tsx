"use client";

import { useState } from "react";
import PostItem from "./PostItem";
import NewPostForm from "./NewPostForm";
import styles from "../page.module.css";

type Post = {
  id: string;
  title: string;
  body: string;
  updatedAt: string;
};

type PostListProps = {
  posts: Post[];
};

export default function PostList({ posts }: PostListProps) {
  const [showNewPostForm, setShowNewPostForm] = useState(false);

  const handleNewPost = () => {
    setShowNewPostForm(true);
  };

  const handleCancelNewPost = () => {
    setShowNewPostForm(false);
  };

  const handleSaveNewPost = () => {
    setShowNewPostForm(false);
  };

  return (
    <>
      <div className={styles.headerActions}>
        <button
          onClick={handleNewPost}
          className={styles.newPostButton}
          disabled={showNewPostForm}
        >
          新しいPostを作成
        </button>
      </div>

      {showNewPostForm && (
        <div className={styles.newPostFormContainer}>
          <NewPostForm
            onCancel={handleCancelNewPost}
            onSave={handleSaveNewPost}
          />
        </div>
      )}

      <div className={styles.postList}>
        {posts.map((item) => (
          <PostItem key={item.id} post={item} />
        ))}
      </div>
    </>
  );
}