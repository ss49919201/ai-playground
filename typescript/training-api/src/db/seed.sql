
INSERT INTO users (id, email, password, name, created_at)
VALUES (
  '1e9a8b7c-6d5e-4f3c-2d1b-0a9b8c7d6e5f',
  'test@example.com',
  '$2a$10$JwZPb5xRHQH6ycnAki7.UuL.5BZ.7SDnS5JGy2JQ.UfvVnKfNn5.q', -- password: password123
  'Test User',
  '2023-01-01T00:00:00.000Z'
);

INSERT INTO training_records (id, title, date, description, created_at, updated_at, user_id)
VALUES (
  'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d',
  '上半身トレーニング',
  '2023-01-15',
  '胸、肩、腕のトレーニング',
  '2023-01-15T10:00:00.000Z',
  '2023-01-15T10:00:00.000Z',
  '1e9a8b7c-6d5e-4f3c-2d1b-0a9b8c7d6e5f'
);

INSERT INTO training_records (id, title, date, description, created_at, updated_at, user_id)
VALUES (
  'b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e',
  '下半身トレーニング',
  '2023-01-17',
  '脚と腹筋のトレーニング',
  '2023-01-17T10:00:00.000Z',
  '2023-01-17T10:00:00.000Z',
  '1e9a8b7c-6d5e-4f3c-2d1b-0a9b8c7d6e5f'
);

INSERT INTO exercises (id, record_id, name)
VALUES (
  'c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f',
  'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d',
  'ベンチプレス'
);

INSERT INTO exercises (id, record_id, name)
VALUES (
  'd4e5f6a7-b8c9-0d1e-2f3a-4b5c6d7e8f9a',
  'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d',
  'ショルダープレス'
);

INSERT INTO exercises (id, record_id, name)
VALUES (
  'e5f6a7b8-c9d0-1e2f-3a4b-5c6d7e8f9a0b',
  'b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e',
  'スクワット'
);

INSERT INTO exercises (id, record_id, name)
VALUES (
  'f6a7b8c9-d0e1-2f3a-4b5c-6d7e8f9a0b1c',
  'b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e',
  'レッグプレス'
);

INSERT INTO sets (id, exercise_id, weight, reps, notes)
VALUES (
  'a7b8c9d0-e1f2-3a4b-5c6d-7e8f9a0b1c2d',
  'c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f',
  80.0,
  10,
  '最初のセット'
);

INSERT INTO sets (id, exercise_id, weight, reps, notes)
VALUES (
  'b8c9d0e1-f2a3-4b5c-6d7e-8f9a0b1c2d3e',
  'c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f',
  85.0,
  8,
  '2セット目'
);

INSERT INTO sets (id, exercise_id, weight, reps, notes)
VALUES (
  'c9d0e1f2-a3b4-5c6d-7e8f-9a0b1c2d3e4f',
  'd4e5f6a7-b8c9-0d1e-2f3a-4b5c6d7e8f9a',
  50.0,
  12,
  '肩のトレーニング'
);

INSERT INTO sets (id, exercise_id, weight, reps, notes)
VALUES (
  'd0e1f2a3-b4c5-6d7e-8f9a-0b1c2d3e4f5a',
  'e5f6a7b8-c9d0-1e2f-3a4b-5c6d7e8f9a0b',
  100.0,
  10,
  '脚のトレーニング'
);

INSERT INTO sets (id, exercise_id, weight, reps, notes)
VALUES (
  'e1f2a3b4-c5d6-7e8f-9a0b-1c2d3e4f5a6b',
  'f6a7b8c9-d0e1-2f3a-4b5c-6d7e8f9a0b1c',
  120.0,
  8,
  'レッグプレス'
);
