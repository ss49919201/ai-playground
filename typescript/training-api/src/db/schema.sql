CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  name TEXT NOT NULL,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS training_records (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  date TEXT NOT NULL,
  description TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  user_id TEXT,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS exercises (
  id TEXT PRIMARY KEY,
  record_id TEXT NOT NULL,
  name TEXT NOT NULL,
  FOREIGN KEY (record_id) REFERENCES training_records(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sets (
  id TEXT PRIMARY KEY,
  exercise_id TEXT NOT NULL,
  weight REAL NOT NULL,
  reps INTEGER NOT NULL,
  notes TEXT,
  FOREIGN KEY (exercise_id) REFERENCES exercises(id) ON DELETE CASCADE
);
