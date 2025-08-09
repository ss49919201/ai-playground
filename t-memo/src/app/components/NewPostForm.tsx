"use client";

import PostForm from "./PostForm";

type NewPostFormProps = {
  onCancel: () => void;
  onSave: () => void;
};

export default function NewPostForm({ onCancel, onSave }: NewPostFormProps) {
  return (
    <PostForm
      onCancel={onCancel}
      onSave={onSave}
    />
  );
}