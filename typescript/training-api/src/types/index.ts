export interface TrainingRecord {
  id: string;
  date: string;
  title: string;
  description: string;
  exercises: Exercise[];
  createdAt: string;
  updatedAt: string;
}

export interface Exercise {
  id: string;
  name: string;
  sets: Set[];
}

export interface Set {
  id: string;
  weight: number; // 重量（kg）
  reps: number; // 回数
  notes?: string; // メモ（任意）
}
