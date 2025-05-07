// ユーザーの型定義
export interface User {
  id: string;
  email: string;
  name?: string;
  createdAt: string;
  updatedAt: string;
}

// セッションの型定義
export interface Session {
  id: string;
  userId: string;
  expiresAt: string;
  createdAt: string;
}

// トレーニング記録の型定義
export interface TrainingRecord {
  id: string;
  userId?: string; // ユーザーID（認証導入後に必須になる予定）
  date: string;
  title: string;
  description: string;
  exercises: Exercise[];
  createdAt: string;
  updatedAt: string;
}

// トレーニング種目の型定義
export interface Exercise {
  id: string;
  name: string;
  sets: Set[];
}

// セットの型定義
export interface Set {
  id: string;
  weight: number; // 重量（kg）
  reps: number; // 回数
  notes?: string; // メモ（任意）
}
