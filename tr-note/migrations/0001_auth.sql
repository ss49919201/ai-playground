-- ユーザーテーブル
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  name TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

-- セッションテーブル
CREATE TABLE sessions (
  id TEXT PRIMARY KEY,
  user_id TEXT NOT NULL,
  expires_at TEXT NOT NULL,
  created_at TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- トレーニング記録テーブルにuser_id外部キーを追加
ALTER TABLE training_records ADD COLUMN user_id TEXT;
ALTER TABLE training_records ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id);