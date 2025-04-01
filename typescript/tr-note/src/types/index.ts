// トレーニング記録の型定義
export interface TrainingRecord {
  id: string;
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
