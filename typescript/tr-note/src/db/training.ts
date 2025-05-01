"use server";

import { Exercise, Set, TrainingRecord } from "../types";
import { getDB, generateId, getCurrentISOString } from "./index";

/**
 * トレーニング記録の所有者かどうかを検証する
 */
export async function verifyTrainingRecordOwnership(
  trainingRecordId: string, 
  userId: string
): Promise<boolean> {
  const db = await getDB();
  
  const record = await db.prepare(
    `SELECT id FROM training_records
     WHERE id = ? AND user_id = ?`
  )
  .bind(trainingRecordId, userId)
  .first<{ id: string }>();
  
  return record !== null;
}

/**
 * トレーニング記録を追加する
 */
export async function addTrainingRecord(
  record: Omit<TrainingRecord, "id" | "createdAt" | "updatedAt">
): Promise<TrainingRecord> {
  const db = await getDB();
  const id = await generateId();
  const timestamp = await getCurrentISOString();
  
  const newRecord: TrainingRecord = {
    ...record,
    id,
    createdAt: timestamp,
    updatedAt: timestamp,
  };
  
  await db.prepare(
    `INSERT INTO training_records (id, user_id, date, title, description, created_at, updated_at) 
     VALUES (?, ?, ?, ?, ?, ?, ?)`
  )
  .bind(
    id,
    record.userId,
    record.date,
    record.title,
    record.description,
    timestamp,
    timestamp
  )
  .run();
  
  // 空の配列を返すだけなので、DBに保存する必要はない
  newRecord.exercises = [];
  
  console.log("トレーニング記録が追加されました:", newRecord);
  return newRecord;
}

/**
 * トレーニング記録を取得する
 */
export async function getTrainingRecord(
  id: string,
  userId?: string
): Promise<TrainingRecord | undefined> {
  const db = await getDB();
  
  // トレーニング記録を取得（ユーザーIDが指定されている場合は所有権チェック）
  let query = `SELECT id, user_id AS userId, date, title, description, created_at AS createdAt, updated_at AS updatedAt
               FROM training_records
               WHERE id = ?`;
  
  const params: any[] = [id];
  
  if (userId) {
    query += ` AND user_id = ?`;
    params.push(userId);
  }
  
  const record = await db.prepare(query)
    .bind(...params)
    .first<Omit<TrainingRecord, "exercises">>();
  
  if (!record) {
    return undefined;
  }
  
  // 関連する種目を取得
  const exercises = await db.prepare(
    `SELECT id, name
     FROM exercises
     WHERE training_record_id = ?`
  )
  .bind(id)
  .all<Omit<Exercise, "sets">>();
  
  // 各種目のセットを取得
  const exercisesWithSets: Exercise[] = [];
  
  for (const exercise of exercises.results) {
    const sets = await db.prepare(
      `SELECT id, weight, reps, notes
       FROM sets
       WHERE exercise_id = ?`
    )
    .bind(exercise.id)
    .all<Set>();
    
    exercisesWithSets.push({
      ...exercise,
      sets: sets.results
    });
  }
  
  // 完全なトレーニング記録を返す
  return {
    ...record,
    exercises: exercisesWithSets
  };
}

/**
 * ユーザーの全てのトレーニング記録を取得する
 */
export async function getAllTrainingRecords(userId: string): Promise<TrainingRecord[]> {
  const db = await getDB();
  
  // 指定ユーザーのトレーニング記録を取得
  const records = await db.prepare(
    `SELECT id, user_id AS userId, date, title, description, created_at AS createdAt, updated_at AS updatedAt
     FROM training_records
     WHERE user_id = ?
     ORDER BY date DESC`
  )
  .bind(userId)
  .all<Omit<TrainingRecord, "exercises">>();
  
  const trainingRecords: TrainingRecord[] = [];
  
  // 各トレーニング記録に対して、種目とセットを取得
  for (const record of records.results) {
    const completeRecord = await getTrainingRecord(record.id, userId);
    if (completeRecord) {
      trainingRecords.push(completeRecord);
    }
  }
  
  return trainingRecords;
}

/**
 * トレーニング記録を削除する
 */
export async function deleteTrainingRecord(id: string, userId: string): Promise<void> {
  const db = await getDB();
  
  // 所有権チェック付きで削除（ユーザーIDが一致する場合のみ削除される）
  // カスケード削除が有効なので、トレーニング記録を削除するだけで関連する種目とセットも削除される
  await db.prepare(
    `DELETE FROM training_records
     WHERE id = ? AND user_id = ?`
  )
  .bind(id, userId)
  .run();
  
  console.log(`ID: ${id} のトレーニング記録が削除されました`);
}

/**
 * トレーニング記録に種目を追加する
 */
export async function addExerciseToRecord(
  recordId: string,
  exercise: Omit<Exercise, "id">,
  userId: string
): Promise<void> {
  const db = await getDB();
  
  // 所有権チェック
  const record = await db.prepare(
    `SELECT id FROM training_records
     WHERE id = ? AND user_id = ?`
  )
  .bind(recordId, userId)
  .first<{ id: string }>();
  
  if (!record) {
    throw new Error(`所有権がないか、トレーニング記録が存在しません。ID: ${recordId}`);
  }
  
  const exerciseId = await generateId();
  const timestamp = await getCurrentISOString();
  
  // 種目を追加
  await db.prepare(
    `INSERT INTO exercises (id, training_record_id, name, created_at)
     VALUES (?, ?, ?, ?)`
  )
  .bind(
    exerciseId,
    recordId,
    exercise.name,
    timestamp
  )
  .run();
  
  // トレーニング記録の更新日時を更新
  await db.prepare(
    `UPDATE training_records
     SET updated_at = ?
     WHERE id = ?`
  )
  .bind(timestamp, recordId)
  .run();
  
  console.log(
    `ID: ${recordId} のトレーニング記録に種目が追加されました:`,
    { id: exerciseId, ...exercise }
  );
}

/**
 * 種目にセットを追加する
 */
export async function addSetToExercise(
  recordId: string,
  exerciseId: string,
  set: Omit<Set, "id">,
  userId: string
): Promise<void> {
  const db = await getDB();
  
  // 所有権チェック
  const record = await db.prepare(
    `SELECT tr.id 
     FROM training_records tr
     JOIN exercises ex ON tr.id = ex.training_record_id
     WHERE tr.id = ? AND tr.user_id = ? AND ex.id = ?`
  )
  .bind(recordId, userId, exerciseId)
  .first<{ id: string }>();
  
  if (!record) {
    throw new Error(`所有権がないか、トレーニング記録または種目が存在しません。Record ID: ${recordId}, Exercise ID: ${exerciseId}`);
  }
  
  const setId = await generateId();
  const timestamp = await getCurrentISOString();
  
  // セットを追加
  await db.prepare(
    `INSERT INTO sets (id, exercise_id, weight, reps, notes, created_at)
     VALUES (?, ?, ?, ?, ?, ?)`
  )
  .bind(
    setId,
    exerciseId,
    set.weight,
    set.reps,
    set.notes || null,
    timestamp
  )
  .run();
  
  // トレーニング記録の更新日時を更新
  await db.prepare(
    `UPDATE training_records
     SET updated_at = ?
     WHERE id = ?`
  )
  .bind(timestamp, recordId)
  .run();
  
  console.log(
    `ID: ${exerciseId} の種目にセットが追加されました:`,
    { id: setId, ...set }
  );
}