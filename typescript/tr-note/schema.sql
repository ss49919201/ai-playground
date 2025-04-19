-- トレーニング記録テーブル
CREATE TABLE training_records (
  id TEXT PRIMARY KEY,
  date TEXT NOT NULL,
  title TEXT NOT NULL,
  description TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

-- トレーニング種目テーブル
CREATE TABLE exercises (
  id TEXT PRIMARY KEY,
  training_record_id TEXT NOT NULL,
  name TEXT NOT NULL,
  created_at TEXT NOT NULL,
  FOREIGN KEY (training_record_id) REFERENCES training_records(id) ON DELETE CASCADE
);

-- セットテーブル
CREATE TABLE sets (
  id TEXT PRIMARY KEY,
  exercise_id TEXT NOT NULL,
  weight REAL NOT NULL,
  reps INTEGER NOT NULL,
  notes TEXT,
  created_at TEXT NOT NULL,
  FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);